package pj

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
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

func writeGraphQLResponse(t *testing.T, w http.ResponseWriter, payload map[string]any) {
	t.Helper()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}
