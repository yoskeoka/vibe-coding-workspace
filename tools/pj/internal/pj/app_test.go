package pj

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInitProvisionsFieldsAndRefreshesCache(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "config.json")
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
		"--config", configPath,
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
	configData, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("read config: %v", err)
	}
	if !strings.Contains(string(configData), `"owner": "yoskeoka"`) {
		t.Fatalf("config contents missing owner: %s", string(configData))
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		t.Fatalf("read cache: %v", err)
	}
	if !strings.Contains(string(data), `"Workspace Repo"`) || !strings.Contains(string(data), `"Priority"`) {
		t.Fatalf("cache contents missing provisioned fields: %s", string(data))
	}
}

func TestRunInitDoesNotWriteStaleCacheWhenProvisioningMutatesAndRefreshFails(t *testing.T) {
	t.Parallel()

	configPath := filepath.Join(t.TempDir(), "config.json")
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
		},
		syncErrs:         []error{nil, os.ErrPermission},
		provisionCreated: true,
		provisionErr:     errStubProvision,
	}

	orig := newProjectClient
	newProjectClient = func() (projectClient, error) { return client, nil }
	defer func() { newProjectClient = orig }()

	err := runInit([]string{
		"--config", configPath,
		"--cache", cachePath,
		"--owner", "yoskeoka",
		"--owner-type", "user",
	}, io.Discard)
	if err == nil {
		t.Fatal("runInit() error = nil, want combined provisioning/refresh error")
	}
	if !strings.Contains(err.Error(), "refusing to write stale cache") {
		t.Fatalf("runInit() error = %q", err)
	}
	if _, statErr := os.Stat(cachePath); !os.IsNotExist(statErr) {
		t.Fatalf("cache file should not exist, stat error = %v", statErr)
	}
}

func TestRunInitRejectsOwnerMismatchWithStoredConfig(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	if err := writeOwnerConfig(configPath, &OwnerConfig{Owner: "yoskeoka", OwnerType: "user"}); err != nil {
		t.Fatalf("writeOwnerConfig(): %v", err)
	}

	err := runInit([]string{
		"--config", configPath,
		"--cache", filepath.Join(dir, "cache.json"),
		"--owner", "acme",
		"--owner-type", "org",
	}, io.Discard)
	if err == nil {
		t.Fatal("runInit() error = nil, want mismatch error")
	}
	if !strings.Contains(err.Error(), "owner config mismatch") {
		t.Fatalf("runInit() error = %q", err)
	}
}

func TestRunSyncUsesStoredOwnerConfigAndCachedProjectNumber(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	cachePath := filepath.Join(dir, "cache.json")
	if err := writeOwnerConfig(configPath, &OwnerConfig{Owner: "yoskeoka", OwnerType: "user"}); err != nil {
		t.Fatalf("writeOwnerConfig(): %v", err)
	}
	if err := writeCache(cachePath, &Cache{
		Project: ProjectRef{
			Owner:         "yoskeoka",
			OwnerType:     "user",
			ProjectNumber: 7,
			ProjectID:     "proj-1",
			Title:         canonicalProjectTitle,
		},
		Fields: compatibleWorkflowFields(),
	}); err != nil {
		t.Fatalf("writeCache(): %v", err)
	}

	client := &stubProjectClient{
		syncResults: []*Cache{
			{
				Project: ProjectRef{
					Owner:         "yoskeoka",
					OwnerType:     "user",
					ProjectNumber: 7,
					ProjectID:     "proj-1",
					Title:         canonicalProjectTitle,
				},
				Fields: compatibleWorkflowFields(),
				Items:  []Item{{ID: "item-1", Title: "Task"}},
			},
		},
	}

	orig := newProjectClient
	newProjectClient = func() (projectClient, error) { return client, nil }
	defer func() { newProjectClient = orig }()

	var stdout strings.Builder
	err := runSync([]string{
		"--config", configPath,
		"--cache", cachePath,
	}, &stdout)
	if err != nil {
		t.Fatalf("runSync() error = %v", err)
	}
	if client.syncCalls != 1 {
		t.Fatalf("sync calls = %d, want 1", client.syncCalls)
	}
	if !strings.Contains(stdout.String(), `synced 1 items from user project "yoskeoka" (#7)`) {
		t.Fatalf("stdout = %q", stdout.String())
	}
}

func TestRunConfigSetClearsMismatchedCache(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	cachePath := filepath.Join(dir, "cache.json")
	if err := writeCache(cachePath, &Cache{
		Project: ProjectRef{
			Owner:         "yoskeoka",
			OwnerType:     "user",
			ProjectNumber: 7,
		},
	}); err != nil {
		t.Fatalf("writeCache(): %v", err)
	}

	err := runConfigSet([]string{
		"--config", configPath,
		"--cache", cachePath,
		"--owner", "acme",
		"--owner-type", "org",
	}, io.Discard)
	if err != nil {
		t.Fatalf("runConfigSet() error = %v", err)
	}
	if _, err := os.Stat(cachePath); !os.IsNotExist(err) {
		t.Fatalf("cache should be removed, stat error = %v", err)
	}
	cfg, err := loadOwnerConfig(configPath)
	if err != nil {
		t.Fatalf("loadOwnerConfig(): %v", err)
	}
	if cfg.Owner != "acme" || cfg.OwnerType != "org" {
		t.Fatalf("config = %+v, want acme/org", cfg)
	}
}

func TestRunConfigClearRemovesConfigAndCache(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.json")
	cachePath := filepath.Join(dir, "cache.json")
	if err := writeOwnerConfig(configPath, &OwnerConfig{Owner: "yoskeoka", OwnerType: "user"}); err != nil {
		t.Fatalf("writeOwnerConfig(): %v", err)
	}
	if err := writeCache(cachePath, &Cache{Project: ProjectRef{Owner: "yoskeoka", OwnerType: "user", ProjectNumber: 7}}); err != nil {
		t.Fatalf("writeCache(): %v", err)
	}

	err := runConfigClear([]string{
		"--config", configPath,
		"--cache", cachePath,
	}, io.Discard)
	if err != nil {
		t.Fatalf("runConfigClear() error = %v", err)
	}
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Fatalf("config should be removed, stat error = %v", err)
	}
	if _, err := os.Stat(cachePath); !os.IsNotExist(err) {
		t.Fatalf("cache should be removed, stat error = %v", err)
	}
}

type stubProjectClient struct {
	ensureProjectRef  ProjectRef
	ensureProjectMade bool
	syncResults       []*Cache
	syncErrs          []error
	syncCalls         int
	provisionCreated  bool
	provisionErr      error
	provisionCalls    int
}

func (s *stubProjectClient) ensureProject(owner, ownerType, title string) (ProjectRef, bool, error) {
	return s.ensureProjectRef, s.ensureProjectMade, nil
}

func (s *stubProjectClient) syncProject(ref ProjectRef) (*Cache, error) {
	var result *Cache
	if s.syncCalls < len(s.syncResults) {
		result = s.syncResults[s.syncCalls]
	}
	var err error
	if s.syncCalls < len(s.syncErrs) {
		err = s.syncErrs[s.syncCalls]
	}
	s.syncCalls++
	return result, err
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

var errStubProvision = os.ErrInvalid
