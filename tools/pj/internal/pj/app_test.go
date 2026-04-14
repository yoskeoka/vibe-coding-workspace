package pj

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInitProvisionsFieldsAndRefreshesCache(t *testing.T) {
	t.Parallel()

	cachePath := filepath.Join(t.TempDir(), "cache.json")
	client := &stubProjectClient{
		ensureProjectRef: ProjectRef{
			Owner:         "yoskeoka",
			OwnerType:     "user",
			ProjectNumber: 7,
			ProjectID:     "proj-1",
			Title:         canonicalProjectTitle,
		},
		syncResults: []*Cache{
			{
				Project: ProjectRef{
					Owner:         "yoskeoka",
					OwnerType:     "user",
					ProjectNumber: 7,
					ProjectID:     "proj-1",
					Title:         canonicalProjectTitle,
				},
				Fields: map[string]FieldCache{
					fieldStatus: {ID: "status", Type: fieldTypeSingleSelect},
				},
			},
			{
				Project: ProjectRef{
					Owner:         "yoskeoka",
					OwnerType:     "user",
					ProjectNumber: 7,
					ProjectID:     "proj-1",
					Title:         canonicalProjectTitle,
				},
				Fields: compatibleWorkflowFields(),
			},
		},
		provisionCreated: true,
	}

	orig := newProjectClient
	newProjectClient = func() (projectClient, error) { return client, nil }
	defer func() { newProjectClient = orig }()

	var stdout strings.Builder
	err := runInit([]string{
		"--cache", cachePath,
		"--owner", "yoskeoka",
		"--owner-type", "user",
	}, &stdout)
	if err != nil {
		t.Fatalf("runInit() error = %v", err)
	}
	if client.syncCalls != 2 {
		t.Fatalf("sync calls = %d, want 2", client.syncCalls)
	}
	if client.provisionCalls != 1 {
		t.Fatalf("provision calls = %d, want 1", client.provisionCalls)
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatalf("read cache: %v", err)
	}
	if !strings.Contains(string(data), `"Workspace Repo"`) || !strings.Contains(string(data), `"Priority"`) {
		t.Fatalf("cache contents missing provisioned fields: %s", string(data))
	}
}

type stubProjectClient struct {
	ensureProjectRef  ProjectRef
	ensureProjectMade bool
	syncResults       []*Cache
	syncCalls         int
	provisionCreated  bool
	provisionErr      error
	provisionCalls    int
}

func (s *stubProjectClient) ensureProject(owner, ownerType, title string) (ProjectRef, bool, error) {
	return s.ensureProjectRef, s.ensureProjectMade, nil
}

func (s *stubProjectClient) syncProject(ref ProjectRef) (*Cache, error) {
	result := s.syncResults[s.syncCalls]
	s.syncCalls++
	return result, nil
}

func (s *stubProjectClient) provisionWorkflowFields(cache *Cache) (bool, error) {
	s.provisionCalls++
	return s.provisionCreated, s.provisionErr
}

func (s *stubProjectClient) addDraftItem(cache *Cache, title, body string, fieldValues map[string]string) error {
	return nil
}

func (s *stubProjectClient) moveItem(cache *Cache, itemID, status string) error {
	return nil
}
