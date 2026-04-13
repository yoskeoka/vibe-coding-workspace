package pj

import (
	"os"
	"path/filepath"
	"time"
)

const (
	defaultCacheRelPath   = ".local/pj/cache.json"
	canonicalProjectTitle = "Workspace Task Triage"
	fieldStatus           = "Status"
	fieldRepo             = "Repo"
	fieldKind             = "Kind"
	fieldPriority         = "Priority"
)

var defaultCachePath = resolveDefaultCachePath()

type ProjectRef struct {
	Owner         string `json:"owner"`
	OwnerType     string `json:"owner_type"`
	ProjectNumber int    `json:"project_number"`
	ProjectID     string `json:"project_id"`
	Title         string `json:"title"`
}

type FieldCache struct {
	ID      string            `json:"id"`
	Type    string            `json:"type"`
	Options map[string]string `json:"options,omitempty"`
}

type Item struct {
	ID          string `json:"id"`
	ContentType string `json:"content_type"`
	Title       string `json:"title"`
	Body        string `json:"body,omitempty"`
	URL         string `json:"url,omitempty"`
	Repository  string `json:"repository,omitempty"`
	Status      string `json:"status,omitempty"`
	Repo        string `json:"repo,omitempty"`
	Kind        string `json:"kind,omitempty"`
	Priority    string `json:"priority,omitempty"`
}

type Cache struct {
	SyncedAt time.Time             `json:"synced_at"`
	Project  ProjectRef            `json:"project"`
	Fields   map[string]FieldCache `json:"fields"`
	Items    []Item                `json:"items"`
}

var requiredFieldNames = []string{
	fieldStatus,
	fieldRepo,
	fieldKind,
	fieldPriority,
}

func resolveDefaultCachePath() string {
	wd, err := os.Getwd()
	if err != nil {
		return defaultCacheRelPath
	}

	root, ok := findGitRoot(wd)
	if !ok {
		return defaultCacheRelPath
	}
	return filepath.Join(root, defaultCacheRelPath)
}

func findGitRoot(start string) (string, bool) {
	dir := start
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, true
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}
