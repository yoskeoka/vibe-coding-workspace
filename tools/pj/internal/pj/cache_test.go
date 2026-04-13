package pj

import (
	"path/filepath"
	"testing"
	"time"
)

func TestWriteAndLoadCache(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "cache.json")
	want := &Cache{
		SyncedAt: time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC),
		Project: ProjectRef{
			Owner:         "yoskeoka",
			OwnerType:     "user",
			ProjectNumber: 7,
			ProjectID:     "PVT_123",
			Title:         "Workspace",
		},
		Fields: map[string]FieldCache{
			fieldStatus: {
				ID:      "field-1",
				Type:    "ProjectV2SingleSelectField",
				Options: map[string]string{"Todo": "opt-1"},
			},
		},
		Items: []Item{{ID: "item-1", Title: "Example", Status: "Todo"}},
	}

	if err := writeCache(path, want); err != nil {
		t.Fatalf("writeCache() error = %v", err)
	}

	got, err := loadCache(path)
	if err != nil {
		t.Fatalf("loadCache() error = %v", err)
	}

	if got.Project.ProjectID != want.Project.ProjectID {
		t.Fatalf("project id = %q, want %q", got.Project.ProjectID, want.Project.ProjectID)
	}
	if got.Items[0].Title != want.Items[0].Title {
		t.Fatalf("item title = %q, want %q", got.Items[0].Title, want.Items[0].Title)
	}
	if got.Fields[fieldStatus].Options["Todo"] != "opt-1" {
		t.Fatalf("status option id = %q, want %q", got.Fields[fieldStatus].Options["Todo"], "opt-1")
	}
}
