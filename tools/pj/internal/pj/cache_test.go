package pj

import (
	"os"
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

func TestFindGitRoot(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git: %v", err)
	}

	nested := filepath.Join(root, "tools", "pj")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatalf("mkdir nested: %v", err)
	}

	got, ok := findGitRoot(nested)
	if !ok {
		t.Fatal("findGitRoot() = not found, want found")
	}
	if got != root {
		t.Fatalf("findGitRoot() = %q, want %q", got, root)
	}
}

func TestMergeProjectRef(t *testing.T) {
	t.Parallel()

	got := mergeProjectRef(ProjectRef{Owner: "yoskeoka"}, ProjectRef{
		Owner:         "ignored",
		OwnerType:     "user",
		ProjectNumber: 7,
	})

	if got.Owner != "yoskeoka" {
		t.Fatalf("owner = %q, want yoskeoka", got.Owner)
	}
	if got.OwnerType != "user" {
		t.Fatalf("owner type = %q, want user", got.OwnerType)
	}
	if got.ProjectNumber != 7 {
		t.Fatalf("project number = %d, want 7", got.ProjectNumber)
	}
}

func TestValidateRequiredFields(t *testing.T) {
	t.Parallel()

	cache := &Cache{
		Fields: map[string]FieldCache{
			fieldStatus: {ID: "status"},
			fieldRepo:   {ID: "repo"},
		},
	}

	err := validateRequiredFields(cache)
	if err == nil {
		t.Fatal("validateRequiredFields() error = nil, want missing-field error")
	}
	got := err.Error()
	if got != "project is missing required fields: Kind, Priority" {
		t.Fatalf("validateRequiredFields() error = %q", got)
	}
}
