package pj

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"
)

func TestEnsureProjectReturnsExistingProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		if !strings.Contains(string(body), "projectsV2(first: 100)") {
			t.Fatalf("unexpected graphql query: %s", string(body))
		}

		writeGraphQLResponse(t, w, map[string]any{
			"data": map[string]any{
				"user": map[string]any{
					"id": "owner-1",
					"projectsV2": map[string]any{
						"pageInfo": map[string]any{"hasNextPage": false},
						"nodes": []map[string]any{
							{"id": "proj-1", "title": canonicalProjectTitle, "number": 42},
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}

	ref, created, err := client.ensureProject("yoskeoka", "user", canonicalProjectTitle)
	if err != nil {
		t.Fatalf("ensureProject() error = %v", err)
	}
	if created {
		t.Fatal("ensureProject() created = true, want false")
	}
	if ref.ProjectID != "proj-1" || ref.ProjectNumber != 42 {
		t.Fatalf("ensureProject() project = %+v", ref)
	}
}

func TestEnsureProjectCreatesMissingProject(t *testing.T) {
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		requests++

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}

		switch requests {
		case 1:
			if !strings.Contains(string(body), "projectsV2(first: 100)") {
				t.Fatalf("unexpected graphql query: %s", string(body))
			}
			writeGraphQLResponse(t, w, map[string]any{
				"data": map[string]any{
					"user": map[string]any{
						"id": "owner-1",
						"projectsV2": map[string]any{
							"pageInfo": map[string]any{"hasNextPage": false},
							"nodes":    []map[string]any{},
						},
					},
				},
			})
		case 2:
			if !strings.Contains(string(body), "createProjectV2") {
				t.Fatalf("unexpected graphql mutation: %s", string(body))
			}
			writeGraphQLResponse(t, w, map[string]any{
				"data": map[string]any{
					"createProjectV2": map[string]any{
						"projectV2": map[string]any{
							"id":     "proj-2",
							"title":  canonicalProjectTitle,
							"number": 43,
						},
					},
				},
			})
		default:
			t.Fatalf("unexpected request count %d", requests)
		}
	}))
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}

	ref, created, err := client.ensureProject("yoskeoka", "user", canonicalProjectTitle)
	if err != nil {
		t.Fatalf("ensureProject() error = %v", err)
	}
	if !created {
		t.Fatal("ensureProject() created = false, want true")
	}
	if ref.ProjectID != "proj-2" || ref.ProjectNumber != 43 {
		t.Fatalf("ensureProject() project = %+v", ref)
	}
}

func TestSyncProjectPaginatesFields(t *testing.T) {
	recorder := newGraphQLTestRecorder(func(w http.ResponseWriter, payload map[string]any, requestNumber int) error {
		switch requestNumber {
		case 1:
			if !strings.Contains(payload["query"].(string), "projectV2(number: $number)") {
				return fmt.Errorf("unexpected project shell query: %s", payload["query"].(string))
			}
			writeGraphQLTestResponse(w, map[string]any{
				"data": map[string]any{
					"user": map[string]any{
						"projectV2": map[string]any{
							"id":    "proj-1",
							"title": canonicalProjectTitle,
						},
					},
				},
			})
		case 2:
			vars := payload["variables"].(map[string]any)
			if vars["after"] != nil {
				return fmt.Errorf("first fields cursor = %#v, want nil", vars["after"])
			}
			writeGraphQLTestResponse(w, map[string]any{
				"data": map[string]any{
					"node": map[string]any{
						"fields": map[string]any{
							"pageInfo": map[string]any{"hasNextPage": true, "endCursor": "fields-1"},
							"nodes": []map[string]any{
								{
									"__typename": "ProjectV2SingleSelectField",
									"id":         "status",
									"name":       fieldStatus,
									"options": []map[string]any{
										{"id": "todo", "name": "Todo"},
									},
								},
							},
						},
					},
				},
			})
		case 3:
			vars := payload["variables"].(map[string]any)
			if vars["after"] != "fields-1" {
				return fmt.Errorf("second fields cursor = %#v, want fields-1", vars["after"])
			}
			writeGraphQLTestResponse(w, map[string]any{
				"data": map[string]any{
					"node": map[string]any{
						"fields": map[string]any{
							"pageInfo": map[string]any{"hasNextPage": false, "endCursor": "fields-2"},
							"nodes": []map[string]any{
								{
									"__typename": "ProjectV2SingleSelectField",
									"id":         "priority",
									"name":       fieldPriority,
									"options": []map[string]any{
										{"id": "high", "name": "High"},
									},
								},
							},
						},
					},
				},
			})
		case 4:
			writeGraphQLTestResponse(w, map[string]any{
				"data": map[string]any{
					"node": map[string]any{
						"items": map[string]any{
							"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
							"nodes":    []map[string]any{},
						},
					},
				},
			})
		default:
			return fmt.Errorf("unexpected request count %d", requestNumber)
		}
		return nil
	})
	server := httptest.NewServer(recorder)
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}

	cache, err := client.syncProject(ProjectRef{Owner: "yoskeoka", OwnerType: "user", ProjectNumber: 42})
	if err != nil {
		t.Fatalf("syncProject() error = %v", err)
	}
	recorder.assertNoError(t)
	requests := recorder.requestsSnapshot()
	if len(requests) != 4 {
		t.Fatalf("request count = %d, want 4", len(requests))
	}
	if cache.Fields[fieldStatus].Options["Todo"] != "todo" {
		t.Fatalf("status field = %+v", cache.Fields[fieldStatus])
	}
	if cache.Fields[fieldPriority].Options["High"] != "high" {
		t.Fatalf("priority field = %+v", cache.Fields[fieldPriority])
	}
}

func TestSyncProjectPaginatesItems(t *testing.T) {
	recorder := newGraphQLTestRecorder(func(w http.ResponseWriter, payload map[string]any, requestNumber int) error {
		switch requestNumber {
		case 1:
			writeGraphQLTestResponse(w, syncProjectShellResponse())
		case 2:
			writeGraphQLTestResponse(w, emptyFieldsResponse())
		case 3:
			vars := payload["variables"].(map[string]any)
			if vars["after"] != nil {
				return fmt.Errorf("first items cursor = %#v, want nil", vars["after"])
			}
			writeGraphQLTestResponse(w, map[string]any{
				"data": map[string]any{
					"node": map[string]any{
						"items": map[string]any{
							"pageInfo": map[string]any{"hasNextPage": true, "endCursor": "items-1"},
							"nodes": []map[string]any{
								syncDraftItemNode("item-1", "First task", []map[string]any{
									syncSingleSelectFieldValue(fieldStatus, "Todo"),
								}, false, nil),
							},
						},
					},
				},
			})
		case 4:
			vars := payload["variables"].(map[string]any)
			if vars["after"] != "items-1" {
				return fmt.Errorf("second items cursor = %#v, want items-1", vars["after"])
			}
			writeGraphQLTestResponse(w, map[string]any{
				"data": map[string]any{
					"node": map[string]any{
						"items": map[string]any{
							"pageInfo": map[string]any{"hasNextPage": false, "endCursor": "items-2"},
							"nodes": []map[string]any{
								syncDraftItemNode("item-2", "Second task", []map[string]any{
									syncSingleSelectFieldValue(fieldPriority, "High"),
								}, false, nil),
							},
						},
					},
				},
			})
		default:
			return fmt.Errorf("unexpected request count %d", requestNumber)
		}
		return nil
	})
	server := httptest.NewServer(recorder)
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}

	cache, err := client.syncProject(ProjectRef{Owner: "yoskeoka", OwnerType: "user", ProjectNumber: 42})
	if err != nil {
		t.Fatalf("syncProject() error = %v", err)
	}
	recorder.assertNoError(t)
	if len(cache.Items) != 2 {
		t.Fatalf("item count = %d, want 2", len(cache.Items))
	}
	if cache.Items[0].Title != "First task" || cache.Items[0].Status != "Todo" {
		t.Fatalf("first item = %+v", cache.Items[0])
	}
	if cache.Items[1].Title != "Second task" || cache.Items[1].Priority != "High" {
		t.Fatalf("second item = %+v", cache.Items[1])
	}
}

func TestSyncProjectPaginatesItemFieldValues(t *testing.T) {
	recorder := newGraphQLTestRecorder(func(w http.ResponseWriter, payload map[string]any, requestNumber int) error {
		switch requestNumber {
		case 1:
			writeGraphQLTestResponse(w, syncProjectShellResponse())
		case 2:
			writeGraphQLTestResponse(w, emptyFieldsResponse())
		case 3:
			writeGraphQLTestResponse(w, map[string]any{
				"data": map[string]any{
					"node": map[string]any{
						"items": map[string]any{
							"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
							"nodes": []map[string]any{
								syncDraftItemNode("item-1", "Task", []map[string]any{
									syncSingleSelectFieldValue(fieldStatus, "Todo"),
								}, true, "values-1"),
							},
						},
					},
				},
			})
		case 4:
			if !strings.Contains(payload["query"].(string), "fieldValues(first: 20, after: $after)") {
				return fmt.Errorf("unexpected field values query: %s", payload["query"].(string))
			}
			vars := payload["variables"].(map[string]any)
			if vars["itemId"] != "item-1" || vars["after"] != "values-1" {
				return fmt.Errorf("field value variables = %#v", vars)
			}
			writeGraphQLTestResponse(w, map[string]any{
				"data": map[string]any{
					"node": map[string]any{
						"fieldValues": map[string]any{
							"pageInfo": map[string]any{"hasNextPage": false, "endCursor": "values-2"},
							"nodes": []map[string]any{
								syncSingleSelectFieldValue(fieldPriority, "High"),
							},
						},
					},
				},
			})
		default:
			return fmt.Errorf("unexpected request count %d", requestNumber)
		}
		return nil
	})
	server := httptest.NewServer(recorder)
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}

	cache, err := client.syncProject(ProjectRef{Owner: "yoskeoka", OwnerType: "user", ProjectNumber: 42})
	if err != nil {
		t.Fatalf("syncProject() error = %v", err)
	}
	recorder.assertNoError(t)
	requests := recorder.requestsSnapshot()
	if len(requests) != 4 {
		t.Fatalf("request count = %d, want 4", len(requests))
	}
	if len(cache.Items) != 1 {
		t.Fatalf("item count = %d, want 1", len(cache.Items))
	}
	if cache.Items[0].Status != "Todo" || cache.Items[0].Priority != "High" {
		t.Fatalf("item = %+v", cache.Items[0])
	}
}

func TestSyncProjectFailsWhenPaginatedResponseHasNoCursor(t *testing.T) {
	tests := []struct {
		name        string
		responses   []map[string]any
		wantErrPart string
	}{
		{
			name: "fields",
			responses: []map[string]any{
				syncProjectShellResponse(),
				paginatedFieldsResponse(true, ""),
			},
			wantErrPart: "project fields response reported another page without an end cursor",
		},
		{
			name: "items",
			responses: []map[string]any{
				syncProjectShellResponse(),
				emptyFieldsResponse(),
				paginatedItemsResponse(true, ""),
			},
			wantErrPart: "project items response reported another page without an end cursor",
		},
		{
			name: "item field values",
			responses: []map[string]any{
				syncProjectShellResponse(),
				emptyFieldsResponse(),
				paginatedItemsResponse(false, nil, syncDraftItemNode("item-1", "Task", []map[string]any{
					syncSingleSelectFieldValue(fieldStatus, "Todo"),
				}, true, "")),
			},
			wantErrPart: "item field values for item-1 response reported another page without an end cursor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := newGraphQLTestRecorder(func(w http.ResponseWriter, _ map[string]any, requestNumber int) error {
				if requestNumber > len(tt.responses) {
					return fmt.Errorf("unexpected request count %d", requestNumber)
				}
				writeGraphQLTestResponse(w, tt.responses[requestNumber-1])
				return nil
			})
			server := httptest.NewServer(recorder)
			defer server.Close()

			client := &githubClient{
				httpClient: server.Client(),
				endpoint:   server.URL,
				token:      "token",
			}

			_, err := client.syncProject(ProjectRef{Owner: "yoskeoka", OwnerType: "user", ProjectNumber: 42})
			recorder.assertNoError(t)
			if err == nil {
				t.Fatal("syncProject() error = nil, want malformed cursor error")
			}
			if !strings.Contains(err.Error(), tt.wantErrPart) {
				t.Fatalf("syncProject() error = %q, want %q", err, tt.wantErrPart)
			}
		})
	}
}

func TestProvisionWorkflowFieldsNoOpForCompatibleBoard(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatalf("unexpected GraphQL request")
	}))
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}

	cache := &Cache{
		Project: ProjectRef{ProjectID: "proj-1"},
		Fields:  compatibleWorkflowFields(),
	}

	created, err := client.provisionWorkflowFields(cache)
	if err != nil {
		t.Fatalf("provisionWorkflowFields() error = %v", err)
	}
	if created {
		t.Fatal("provisionWorkflowFields() created = true, want false")
	}
}

func TestProvisionWorkflowFieldsCreatesMissingFields(t *testing.T) {
	t.Parallel()

	var requests []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		requests = append(requests, payload)

		vars := payload["variables"].(map[string]any)
		name := vars["name"].(string)
		writeGraphQLResponse(t, w, map[string]any{
			"data": map[string]any{
				"createProjectV2Field": map[string]any{
					"projectV2Field": map[string]any{
						"__typename": "ProjectV2SingleSelectField",
						"id":         "field-" + name,
						"name":       name,
					},
				},
			},
		})
	}))
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}

	cache := &Cache{
		Project: ProjectRef{ProjectID: "proj-1"},
		Fields: map[string]FieldCache{
			fieldStatus: {ID: "status", Type: fieldTypeSingleSelect},
		},
	}

	created, err := client.provisionWorkflowFields(cache)
	if err != nil {
		t.Fatalf("provisionWorkflowFields() error = %v", err)
	}
	if !created {
		t.Fatal("provisionWorkflowFields() created = false, want true")
	}
	if len(requests) != 3 {
		t.Fatalf("request count = %d, want 3", len(requests))
	}

	var gotNames []string
	for _, payload := range requests {
		query := payload["query"].(string)
		if !strings.Contains(query, "createProjectV2Field") {
			t.Fatalf("unexpected graphql query: %s", query)
		}

		vars := payload["variables"].(map[string]any)
		gotNames = append(gotNames, vars["name"].(string))

		schema := workflowFieldSchemaByName[vars["name"].(string)]
		var gotOptions []string
		for _, raw := range vars["options"].([]any) {
			option := raw.(map[string]any)
			gotOptions = append(gotOptions, option["name"].(string))
		}
		if !slices.Equal(gotOptions, optionNames(schema.Options)) {
			t.Fatalf("options for %q = %v, want %v", schema.Name, gotOptions, optionNames(schema.Options))
		}
	}
	if !slices.Equal(gotNames, []string{fieldRepo, fieldKind, fieldPriority}) {
		t.Fatalf("created field names = %v", gotNames)
	}
}

func TestProvisionWorkflowFieldsCreatesOnlyMissingFields(t *testing.T) {
	t.Parallel()

	var gotNames []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload struct {
			Variables map[string]any `json:"variables"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		gotNames = append(gotNames, payload.Variables["name"].(string))

		writeGraphQLResponse(t, w, map[string]any{
			"data": map[string]any{
				"createProjectV2Field": map[string]any{
					"projectV2Field": map[string]any{
						"__typename": "ProjectV2SingleSelectField",
						"id":         "field",
						"name":       payload.Variables["name"].(string),
					},
				},
			},
		})
	}))
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}

	fields := compatibleWorkflowFields()
	delete(fields, fieldKind)
	delete(fields, fieldPriority)

	created, err := client.provisionWorkflowFields(&Cache{
		Project: ProjectRef{ProjectID: "proj-1"},
		Fields:  fields,
	})
	if err != nil {
		t.Fatalf("provisionWorkflowFields() error = %v", err)
	}
	if !created {
		t.Fatal("provisionWorkflowFields() created = false, want true")
	}
	if !slices.Equal(gotNames, []string{fieldKind, fieldPriority}) {
		t.Fatalf("created field names = %v", gotNames)
	}
}

func TestProvisionWorkflowFieldsFailsForIncompatibleField(t *testing.T) {
	t.Parallel()

	client := &githubClient{
		httpClient: http.DefaultClient,
		endpoint:   "http://example.invalid",
		token:      "token",
	}

	fields := compatibleWorkflowFields()
	fields[fieldRepo] = FieldCache{ID: "repo", Type: "ProjectV2Field"}

	_, err := client.provisionWorkflowFields(&Cache{
		Project: ProjectRef{ProjectID: "proj-1"},
		Fields:  fields,
	})
	if err == nil {
		t.Fatal("provisionWorkflowFields() error = nil, want incompatible-field error")
	}
	if !strings.Contains(err.Error(), `field "Workspace Repo" is incompatible`) || !strings.Contains(err.Error(), `unsupported type`) {
		t.Fatalf("provisionWorkflowFields() error = %q", err)
	}
}

func TestProvisionWorkflowFieldsFailsForMissingRequiredOption(t *testing.T) {
	t.Parallel()

	client := &githubClient{
		httpClient: http.DefaultClient,
		endpoint:   "http://example.invalid",
		token:      "token",
	}

	fields := compatibleWorkflowFields()
	fields[fieldPriority] = FieldCache{
		ID:      "priority",
		Type:    fieldTypeSingleSelect,
		Options: map[string]string{"High": "high"},
	}

	_, err := client.provisionWorkflowFields(&Cache{
		Project: ProjectRef{ProjectID: "proj-1"},
		Fields:  fields,
	})
	if err == nil {
		t.Fatal("provisionWorkflowFields() error = nil, want incompatible-field error")
	}
	if !strings.Contains(err.Error(), `field "Priority" is incompatible`) || !strings.Contains(err.Error(), `missing required options: Medium, Low`) {
		t.Fatalf("provisionWorkflowFields() error = %q", err)
	}
}

func TestUpdateItemUpdatesDraftIssueAndFields(t *testing.T) {
	var requests []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		requests = append(requests, payload)
		writeGraphQLResponse(t, w, map[string]any{
			"data": map[string]any{
				"ok": true,
			},
		})
	}))
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}
	cache := &Cache{
		Project: ProjectRef{ProjectID: "proj-1"},
		Fields: map[string]FieldCache{
			fieldStatus: {
				ID:      "status",
				Type:    fieldTypeSingleSelect,
				Options: map[string]string{"Done": "done-id"},
			},
		},
		Items: []Item{{ID: "item-1", DraftIssueID: "draft-1", ContentType: "DraftIssue"}},
	}

	err := client.updateItem(cache, "item-1", itemUpdate{
		Title:         "New title",
		TitleProvided: true,
		Body:          "New body",
		BodyProvided:  true,
		FieldValues:   map[string]string{fieldStatus: "Done"},
	})
	if err != nil {
		t.Fatalf("updateItem() error = %v", err)
	}
	if len(requests) != 2 {
		t.Fatalf("request count = %d, want 2", len(requests))
	}
	if !strings.Contains(requests[0]["query"].(string), "updateProjectV2DraftIssue") {
		t.Fatalf("first query = %s", requests[0]["query"].(string))
	}
	input := requests[0]["variables"].(map[string]any)["input"].(map[string]any)
	if input["draftIssueId"] != "draft-1" || input["title"] != "New title" || input["body"] != "New body" {
		t.Fatalf("draft update input = %#v", input)
	}
	if !strings.Contains(requests[1]["query"].(string), "updateProjectV2ItemFieldValue") {
		t.Fatalf("second query = %s", requests[1]["query"].(string))
	}
	vars := requests[1]["variables"].(map[string]any)
	if vars["itemId"] != "item-1" || vars["optionId"] != "done-id" {
		t.Fatalf("field update variables = %#v", vars)
	}
}

func TestProjectLinkedRepositoriesQueriesProjectRepositories(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		if !strings.Contains(string(body), "repositories(first: 100)") {
			t.Fatalf("unexpected graphql query: %s", string(body))
		}
		writeGraphQLResponse(t, w, map[string]any{
			"data": map[string]any{
				"node": map[string]any{
					"__typename": "ProjectV2",
					"repositories": map[string]any{
						"pageInfo": map[string]any{"hasNextPage": false},
						"nodes": []map[string]any{
							{"id": "repo-1", "nameWithOwner": "yoskeoka/vibe-coding-workspace"},
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}
	repos, err := client.projectLinkedRepositories("proj-1")
	if err != nil {
		t.Fatalf("projectLinkedRepositories() error = %v", err)
	}
	if !slices.Equal(repos, []RepositoryRef{{ID: "repo-1", NameWithOwner: "yoskeoka/vibe-coding-workspace"}}) {
		t.Fatalf("repos = %+v", repos)
	}
}

func TestResolveRepositoryQueriesRepository(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if !strings.Contains(payload["query"].(string), "repository(owner: $owner, name: $name)") {
			t.Fatalf("unexpected graphql query: %s", payload["query"].(string))
		}
		vars := payload["variables"].(map[string]any)
		if vars["owner"] != "yoskeoka" || vars["name"] != "vibe-coding-workspace" {
			t.Fatalf("variables = %#v", vars)
		}
		writeGraphQLResponse(t, w, map[string]any{
			"data": map[string]any{
				"repository": map[string]any{
					"id":            "repo-1",
					"nameWithOwner": "yoskeoka/vibe-coding-workspace",
				},
			},
		})
	}))
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}
	repo, err := client.resolveRepository("yoskeoka", "vibe-coding-workspace")
	if err != nil {
		t.Fatalf("resolveRepository() error = %v", err)
	}
	if repo != (RepositoryRef{ID: "repo-1", NameWithOwner: "yoskeoka/vibe-coding-workspace"}) {
		t.Fatalf("repo = %+v", repo)
	}
}

func TestResolveRepositoryFailsOnIncompleteRepository(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeGraphQLResponse(t, w, map[string]any{
			"data": map[string]any{
				"repository": map[string]any{
					"id": "repo-1",
				},
			},
		})
	}))
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}
	_, err := client.resolveRepository("yoskeoka", "vibe-coding-workspace")
	if err == nil {
		t.Fatal("resolveRepository() error = nil, want incomplete repository error")
	}
	if !strings.Contains(err.Error(), "repository yoskeoka/vibe-coding-workspace not found") {
		t.Fatalf("error = %q", err)
	}
}

func TestProjectLinkedRepositoriesFailsOnIncompleteRepositoryNode(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeGraphQLResponse(t, w, map[string]any{
			"data": map[string]any{
				"node": map[string]any{
					"__typename": "ProjectV2",
					"repositories": map[string]any{
						"pageInfo": map[string]any{"hasNextPage": false},
						"nodes": []map[string]any{
							{"id": "repo-1"},
						},
					},
				},
			},
		})
	}))
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}
	_, err := client.projectLinkedRepositories("proj-1")
	if err == nil {
		t.Fatal("projectLinkedRepositories() error = nil, want incomplete node error")
	}
	if !strings.Contains(err.Error(), "linked repository node 0 is incomplete") {
		t.Fatalf("error = %q", err)
	}
}

func TestLinkAndUnlinkProjectToRepositoryUseProjectV2Mutations(t *testing.T) {
	t.Parallel()

	var queries []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		queries = append(queries, payload["query"].(string))
		vars := payload["variables"].(map[string]any)
		if vars["projectId"] != "proj-1" || vars["repositoryId"] != "repo-1" {
			t.Fatalf("variables = %#v", vars)
		}
		writeGraphQLResponse(t, w, map[string]any{"data": map[string]any{}})
	}))
	defer server.Close()

	client := &githubClient{
		httpClient: server.Client(),
		endpoint:   server.URL,
		token:      "token",
	}
	repo := RepositoryRef{ID: "repo-1", NameWithOwner: "yoskeoka/vibe-coding-workspace"}
	if err := client.linkProjectToRepository("proj-1", repo); err != nil {
		t.Fatalf("linkProjectToRepository() error = %v", err)
	}
	if err := client.unlinkProjectFromRepository("proj-1", repo); err != nil {
		t.Fatalf("unlinkProjectFromRepository() error = %v", err)
	}
	if len(queries) != 2 {
		t.Fatalf("query count = %d, want 2", len(queries))
	}
	if !strings.Contains(queries[0], "linkProjectV2ToRepository") {
		t.Fatalf("link query = %s", queries[0])
	}
	if !strings.Contains(queries[1], "unlinkProjectV2FromRepository") {
		t.Fatalf("unlink query = %s", queries[1])
	}
}

func compatibleWorkflowFields() map[string]FieldCache {
	return map[string]FieldCache{
		fieldStatus: {
			ID:   "status",
			Type: fieldTypeSingleSelect,
		},
		fieldRepo: {
			ID:      "repo",
			Type:    fieldTypeSingleSelect,
			Options: optionIDs(workflowFieldSchemaByName[fieldRepo].Options),
		},
		fieldKind: {
			ID:      "kind",
			Type:    fieldTypeSingleSelect,
			Options: optionIDs(workflowFieldSchemaByName[fieldKind].Options),
		},
		fieldPriority: {
			ID:      "priority",
			Type:    fieldTypeSingleSelect,
			Options: optionIDs(workflowFieldSchemaByName[fieldPriority].Options),
		},
	}
}

func optionIDs(options []workflowFieldOption) map[string]string {
	ids := make(map[string]string, len(options))
	for _, option := range options {
		ids[option.Name] = strings.ToLower(option.Name)
	}
	return ids
}

func optionNames(options []workflowFieldOption) []string {
	names := make([]string, 0, len(options))
	for _, option := range options {
		names = append(names, option.Name)
	}
	return names
}

func syncProjectShellResponse() map[string]any {
	return map[string]any{
		"data": map[string]any{
			"user": map[string]any{
				"projectV2": map[string]any{
					"id":    "proj-1",
					"title": canonicalProjectTitle,
				},
			},
		},
	}
}

func emptyFieldsResponse() map[string]any {
	return map[string]any{
		"data": map[string]any{
			"node": map[string]any{
				"fields": map[string]any{
					"pageInfo": map[string]any{"hasNextPage": false, "endCursor": nil},
					"nodes":    []map[string]any{},
				},
			},
		},
	}
}

func paginatedFieldsResponse(hasNextPage bool, endCursor any) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"node": map[string]any{
				"fields": map[string]any{
					"pageInfo": map[string]any{"hasNextPage": hasNextPage, "endCursor": endCursor},
					"nodes":    []map[string]any{},
				},
			},
		},
	}
}

func paginatedItemsResponse(hasNextPage bool, endCursor any, nodes ...map[string]any) map[string]any {
	return map[string]any{
		"data": map[string]any{
			"node": map[string]any{
				"items": map[string]any{
					"pageInfo": map[string]any{"hasNextPage": hasNextPage, "endCursor": endCursor},
					"nodes":    nodes,
				},
			},
		},
	}
}

func syncDraftItemNode(itemID, title string, values []map[string]any, hasMoreValues bool, valuesCursor any) map[string]any {
	return map[string]any{
		"id": itemID,
		"content": map[string]any{
			"__typename": "DraftIssue",
			"id":         "draft-" + itemID,
			"title":      title,
			"body":       "",
		},
		"fieldValues": map[string]any{
			"pageInfo": map[string]any{"hasNextPage": hasMoreValues, "endCursor": valuesCursor},
			"nodes":    values,
		},
	}
}

func syncSingleSelectFieldValue(fieldName, value string) map[string]any {
	return map[string]any{
		"__typename": "ProjectV2ItemFieldSingleSelectValue",
		"name":       value,
		"field": map[string]any{
			"name": fieldName,
		},
	}
}

type graphQLTestRecorder struct {
	mu       sync.Mutex
	requests []map[string]any
	err      error
	handler  func(http.ResponseWriter, map[string]any, int) error
}

func newGraphQLTestRecorder(handler func(http.ResponseWriter, map[string]any, int) error) *graphQLTestRecorder {
	return &graphQLTestRecorder{handler: handler}
}

func (r *graphQLTestRecorder) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var payload map[string]any
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		r.recordError(fmt.Errorf("decode request body: %w", err))
		http.Error(w, "decode request body", http.StatusBadRequest)
		return
	}

	r.mu.Lock()
	r.requests = append(r.requests, payload)
	requestNumber := len(r.requests)
	r.mu.Unlock()

	if err := r.handler(w, payload, requestNumber); err != nil {
		r.recordError(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *graphQLTestRecorder) recordError(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err == nil {
		r.err = err
	}
}

func (r *graphQLTestRecorder) assertNoError(t *testing.T) {
	t.Helper()

	r.mu.Lock()
	defer r.mu.Unlock()
	if r.err != nil {
		t.Fatal(r.err)
	}
}

func (r *graphQLTestRecorder) requestsSnapshot() []map[string]any {
	r.mu.Lock()
	defer r.mu.Unlock()

	requests := make([]map[string]any, len(r.requests))
	copy(requests, r.requests)
	return requests
}

func writeGraphQLTestResponse(w http.ResponseWriter, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}

func writeGraphQLResponse(t *testing.T, w http.ResponseWriter, payload map[string]any) {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}
