package pj

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInitProvisionsFieldsAndRefreshesCache(t *testing.T) {
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

func TestLoadOwnerConfigRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte("{\n  \"owner\": \"yoskeoka\"\n}\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(): %v", err)
	}

	_, err := loadOwnerConfig(path)
	if err == nil {
		t.Fatal("loadOwnerConfig() error = nil, want invalid config error")
	}
	if !strings.Contains(err.Error(), "invalid owner config") {
		t.Fatalf("loadOwnerConfig() error = %q", err)
	}
}

func TestRunConfigShowRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte("{\n  \"owner\": \"yoskeoka\",\n  \"owner_type\": \"\"\n}\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(): %v", err)
	}

	err := runConfigShow([]string{"--config", path}, io.Discard)
	if err == nil {
		t.Fatal("runConfigShow() error = nil, want invalid config error")
	}
	if !strings.Contains(err.Error(), "invalid owner config") {
		t.Fatalf("runConfigShow() error = %q", err)
	}
}

func TestRunAddReadsBodyFileAndResolvesFieldValues(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "cache.json")
	bodyPath := filepath.Join(dir, "body.md")
	if err := os.WriteFile(bodyPath, []byte("handoff body\n"), 0o644); err != nil {
		t.Fatalf("WriteFile(): %v", err)
	}
	if err := writeCache(cachePath, testProjectCache()); err != nil {
		t.Fatalf("writeCache(): %v", err)
	}

	client := &stubProjectClient{
		syncResults: []*Cache{testProjectCache()},
	}
	orig := newProjectClient
	newProjectClient = func() (projectClient, error) { return client, nil }
	defer func() { newProjectClient = orig }()

	err := runAdd([]string{
		"--cache", cachePath,
		"--title", "Task",
		"--body-file", bodyPath,
		"--status", "in-progress",
		"--repo", "ww",
		"--kind", "feature",
		"--priority", "hi",
	}, io.Discard)
	if err != nil {
		t.Fatalf("runAdd() error = %v", err)
	}
	if client.addBody != "handoff body\n" {
		t.Fatalf("add body = %q", client.addBody)
	}
	want := map[string]string{
		fieldStatus:   "In Progress",
		fieldRepo:     "ww",
		fieldKind:     "Feature",
		fieldPriority: "High",
	}
	for field, value := range want {
		if client.addFieldValues[field] != value {
			t.Fatalf("field %q = %q, want %q", field, client.addFieldValues[field], value)
		}
	}
}

func TestRunAddRejectsBodyAndBodyFileTogether(t *testing.T) {
	err := runAdd([]string{
		"--cache", filepath.Join(t.TempDir(), "cache.json"),
		"--title", "Task",
		"--body", "inline",
		"--body-file", "body.md",
	}, io.Discard)
	if err == nil {
		t.Fatal("runAdd() error = nil, want conflict error")
	}
	if !strings.Contains(err.Error(), "either --body or --body-file") {
		t.Fatalf("runAdd() error = %q", err)
	}
}

func TestRunUpdatePartialFieldsAndRepoIndex(t *testing.T) {
	dir := t.TempDir()
	cachePath := filepath.Join(dir, "cache.json")
	bodyPath := filepath.Join(dir, "body.md")
	if err := os.WriteFile(bodyPath, []byte("updated body"), 0o644); err != nil {
		t.Fatalf("WriteFile(): %v", err)
	}
	if err := writeCache(cachePath, testProjectCache()); err != nil {
		t.Fatalf("writeCache(): %v", err)
	}

	client := &stubProjectClient{
		syncResults: []*Cache{testProjectCache()},
	}
	orig := newProjectClient
	newProjectClient = func() (projectClient, error) { return client, nil }
	defer func() { newProjectClient = orig }()

	err := runUpdate([]string{
		"--cache", cachePath,
		"--item", "item-1",
		"--title", "New title",
		"--body-file", bodyPath,
		"--status", "todo",
		"--repo", "2",
	}, io.Discard)
	if err != nil {
		t.Fatalf("runUpdate() error = %v", err)
	}
	if client.updateItemID != "item-1" {
		t.Fatalf("update item = %q", client.updateItemID)
	}
	if !client.update.TitleProvided || client.update.Title != "New title" {
		t.Fatalf("update title = provided %v value %q", client.update.TitleProvided, client.update.Title)
	}
	if !client.update.BodyProvided || client.update.Body != "updated body" {
		t.Fatalf("update body = provided %v value %q", client.update.BodyProvided, client.update.Body)
	}
	if client.update.FieldValues[fieldStatus] != "Todo" {
		t.Fatalf("status = %q", client.update.FieldValues[fieldStatus])
	}
	if client.update.FieldValues[fieldRepo] != "ww" {
		t.Fatalf("repo = %q, want ww", client.update.FieldValues[fieldRepo])
	}
}

func TestRunUpdateRejectsAmbiguousPrefix(t *testing.T) {
	cache := testProjectCache()
	cache.RepoOptions = append(cache.RepoOptions, RepoOption{
		DisplayValue:  "ww-tools",
		SourceType:    repoSourceGitHub,
		SourceURL:     "https://github.com/yoskeoka/ww-tools",
		CanonicalSlug: "github.com/yoskeoka/ww-tools",
		Aliases:       []string{"ww-tools", "yoskeoka/ww-tools"},
	})
	cache.Fields[fieldRepo].Options["ww-tools"] = "ww-tools"

	cachePath := filepath.Join(t.TempDir(), "cache.json")
	if err := writeCache(cachePath, cache); err != nil {
		t.Fatalf("writeCache(): %v", err)
	}

	err := runUpdate([]string{"--cache", cachePath, "--item", "item-1", "--repo", "github.com/yoskeoka/w"}, io.Discard)
	if err == nil {
		t.Fatal("runUpdate() error = nil, want ambiguity error")
	}
	if !strings.Contains(err.Error(), "ambiguous Workspace Repo value") || !strings.Contains(err.Error(), "ww-tools") {
		t.Fatalf("runUpdate() error = %q", err)
	}
}

func TestResolveRepoValueSupportsSlugAliasAndPrefix(t *testing.T) {
	cache := testProjectCache()
	cases := map[string]string{
		"github.com/yoskeoka/ww": "ww",
		"yoskeoka/ww":            "ww",
		"ai_ar":                  "ai-arena",
	}
	for input, want := range cases {
		got, err := resolveRepoValue(cache, input)
		if err != nil {
			t.Fatalf("resolveRepoValue(%q) error = %v", input, err)
		}
		if got != want {
			t.Fatalf("resolveRepoValue(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestRunURLAndOpenUseCanonicalProjectURL(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "cache.json")
	if err := writeCache(cachePath, testProjectCache()); err != nil {
		t.Fatalf("writeCache(): %v", err)
	}

	var stdout strings.Builder
	if err := runURL([]string{"--cache", cachePath}, &stdout); err != nil {
		t.Fatalf("runURL() error = %v", err)
	}
	wantURL := "https://github.com/users/yoskeoka/projects/7"
	if strings.TrimSpace(stdout.String()) != wantURL {
		t.Fatalf("url output = %q, want %q", strings.TrimSpace(stdout.String()), wantURL)
	}

	orig := openProjectURL
	var opened string
	openProjectURL = func(url string) error {
		opened = url
		return nil
	}
	defer func() { openProjectURL = orig }()

	stdout.Reset()
	if err := runOpen([]string{"--cache", cachePath}, &stdout); err != nil {
		t.Fatalf("runOpen() error = %v", err)
	}
	if opened != wantURL {
		t.Fatalf("opened URL = %q, want %q", opened, wantURL)
	}
}

func TestRunRepoLinkStatusReportsLinkedRepository(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "cache.json")
	if err := writeCache(cachePath, testProjectCache()); err != nil {
		t.Fatalf("writeCache(): %v", err)
	}

	client := &stubProjectClient{
		resolvedRepo: RepositoryRef{ID: "repo-1", NameWithOwner: "yoskeoka/vibe-coding-workspace"},
		linkedRepos:  []RepositoryRef{{ID: "repo-1", NameWithOwner: "yoskeoka/vibe-coding-workspace"}},
	}
	orig := newProjectClient
	newProjectClient = func() (projectClient, error) { return client, nil }
	defer func() { newProjectClient = orig }()

	var stdout strings.Builder
	err := runRepoLink([]string{"status", "--cache", cachePath, "yoskeoka/vibe-coding-workspace"}, &stdout)
	if err != nil {
		t.Fatalf("runRepoLink() error = %v", err)
	}
	if !strings.Contains(stdout.String(), "is linked") {
		t.Fatalf("stdout = %q, want linked status", stdout.String())
	}
}

func TestRunRepoLinkAddLinksUnlinkedRepository(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "cache.json")
	if err := writeCache(cachePath, testProjectCache()); err != nil {
		t.Fatalf("writeCache(): %v", err)
	}

	client := &stubProjectClient{
		resolvedRepo: RepositoryRef{ID: "repo-1", NameWithOwner: "yoskeoka/vibe-coding-workspace"},
	}
	orig := newProjectClient
	newProjectClient = func() (projectClient, error) { return client, nil }
	defer func() { newProjectClient = orig }()

	var stdout strings.Builder
	err := runRepoLink([]string{"add", "--cache", cachePath, "yoskeoka/vibe-coding-workspace"}, &stdout)
	if err != nil {
		t.Fatalf("runRepoLink() error = %v", err)
	}
	if client.linkedRepo.ID != "repo-1" {
		t.Fatalf("linked repo = %+v, want repo-1", client.linkedRepo)
	}
	if !strings.Contains(stdout.String(), "linked yoskeoka/vibe-coding-workspace") {
		t.Fatalf("stdout = %q, want linked message", stdout.String())
	}
}

func TestRunRepoLinkRejectsDifferentRepositoryOwner(t *testing.T) {
	cachePath := filepath.Join(t.TempDir(), "cache.json")
	if err := writeCache(cachePath, testProjectCache()); err != nil {
		t.Fatalf("writeCache(): %v", err)
	}

	err := runRepoLink([]string{"status", "--cache", cachePath, "acme/vibe-coding-workspace"}, io.Discard)
	if err == nil {
		t.Fatal("runRepoLink() error = nil, want owner mismatch")
	}
	if !strings.Contains(err.Error(), "does not match project owner") {
		t.Fatalf("error = %q, want owner mismatch", err)
	}
}

func testProjectCache() *Cache {
	return &Cache{
		Project: ProjectRef{
			Owner:         "yoskeoka",
			OwnerType:     "user",
			ProjectNumber: 7,
			ProjectID:     "proj-1",
			Title:         canonicalProjectTitle,
		},
		Fields: map[string]FieldCache{
			fieldStatus: {
				ID:      "status",
				Type:    fieldTypeSingleSelect,
				Options: map[string]string{"Todo": "todo", "In Progress": "in-progress", "Done": "done"},
			},
			fieldRepo: {
				ID:      "repo",
				Type:    fieldTypeSingleSelect,
				Options: map[string]string{"ai-arena": "ai-arena", "ww": "ww"},
			},
			fieldKind: {
				ID:      "kind",
				Type:    fieldTypeSingleSelect,
				Options: map[string]string{"Feature": "feature", "Bug": "bug"},
			},
			fieldPriority: {
				ID:      "priority",
				Type:    fieldTypeSingleSelect,
				Options: map[string]string{"High": "high", "Medium": "medium"},
			},
		},
		RepoOptions: []RepoOption{
			{
				DisplayValue:  "ai-arena",
				SourceType:    repoSourceGitHub,
				SourceURL:     "https://github.com/yoskeoka/ai-arena",
				CanonicalSlug: "github.com/yoskeoka/ai-arena",
				Aliases:       []string{"ai-arena", "yoskeoka/ai-arena"},
			},
			{
				DisplayValue:  "ww",
				SourceType:    repoSourceGitHub,
				SourceURL:     "https://github.com/yoskeoka/ww",
				CanonicalSlug: "github.com/yoskeoka/ww",
				Aliases:       []string{"ww", "yoskeoka/ww"},
			},
		},
		Items: []Item{
			{ID: "item-1", DraftIssueID: "draft-1", ContentType: "DraftIssue", Title: "Task"},
		},
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
	addBody           string
	addFieldValues    map[string]string
	updateItemID      string
	update            itemUpdate
	resolvedRepo      RepositoryRef
	linkedRepos       []RepositoryRef
	linkedRepo        RepositoryRef
	unlinkedRepo      RepositoryRef
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

func (s *stubProjectClient) projectLinkedRepositories(projectID string) ([]RepositoryRef, error) {
	return s.linkedRepos, nil
}

func (s *stubProjectClient) linkProjectToRepository(projectID string, repo RepositoryRef) error {
	s.linkedRepo = repo
	return nil
}

func (s *stubProjectClient) unlinkProjectFromRepository(projectID string, repo RepositoryRef) error {
	s.unlinkedRepo = repo
	return nil
}

func (s *stubProjectClient) resolveRepository(owner, name string) (RepositoryRef, error) {
	return s.resolvedRepo, nil
}

func (s *stubProjectClient) addDraftItem(cache *Cache, title, body string, fieldValues map[string]string) error {
	s.addBody = body
	s.addFieldValues = fieldValues
	return nil
}

func (s *stubProjectClient) updateItem(cache *Cache, itemID string, update itemUpdate) error {
	s.updateItemID = itemID
	s.update = update
	return nil
}

var errStubProvision = os.ErrInvalid
