package pj

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
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
	var requests []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		requests = append(requests, payload)

		switch len(requests) {
		case 1:
			if !strings.Contains(payload["query"].(string), "projectV2(number: $number)") {
				t.Fatalf("unexpected project shell query: %s", payload["query"].(string))
			}
			writeGraphQLResponse(t, w, map[string]any{
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
				t.Fatalf("first fields cursor = %#v, want nil", vars["after"])
			}
			writeGraphQLResponse(t, w, map[string]any{
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
				t.Fatalf("second fields cursor = %#v, want fields-1", vars["after"])
			}
			writeGraphQLResponse(t, w, map[string]any{
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
			writeGraphQLResponse(t, w, map[string]any{
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
			t.Fatalf("unexpected request count %d", len(requests))
		}
	}))
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
	var requests []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		requests = append(requests, payload)

		switch len(requests) {
		case 1:
			writeGraphQLResponse(t, w, syncProjectShellResponse())
		case 2:
			writeGraphQLResponse(t, w, emptyFieldsResponse())
		case 3:
			vars := payload["variables"].(map[string]any)
			if vars["after"] != nil {
				t.Fatalf("first items cursor = %#v, want nil", vars["after"])
			}
			writeGraphQLResponse(t, w, map[string]any{
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
				t.Fatalf("second items cursor = %#v, want items-1", vars["after"])
			}
			writeGraphQLResponse(t, w, map[string]any{
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
			t.Fatalf("unexpected request count %d", len(requests))
		}
	}))
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
	var requests []map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		requests = append(requests, payload)

		switch len(requests) {
		case 1:
			writeGraphQLResponse(t, w, syncProjectShellResponse())
		case 2:
			writeGraphQLResponse(t, w, emptyFieldsResponse())
		case 3:
			writeGraphQLResponse(t, w, map[string]any{
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
				t.Fatalf("unexpected field values query: %s", payload["query"].(string))
			}
			vars := payload["variables"].(map[string]any)
			if vars["itemId"] != "item-1" || vars["after"] != "values-1" {
				t.Fatalf("field value variables = %#v", vars)
			}
			writeGraphQLResponse(t, w, map[string]any{
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
			t.Fatalf("unexpected request count %d", len(requests))
		}
	}))
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

func writeGraphQLResponse(t *testing.T, w http.ResponseWriter, payload map[string]any) {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}
