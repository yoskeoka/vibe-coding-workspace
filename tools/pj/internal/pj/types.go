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
	fieldRepo             = "Workspace Repo"
	fieldKind             = "Kind"
	fieldPriority         = "Priority"
	fieldTypeSingleSelect = "ProjectV2SingleSelectField"
	fieldDataTypeSingle   = "SINGLE_SELECT"
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

type workflowFieldOption struct {
	Name        string
	Description string
	Color       string
}

type workflowFieldSchema struct {
	Name      string
	Type      string
	Provision bool
	Options   []workflowFieldOption
}

var workflowFieldSchemas = []workflowFieldSchema{
	{
		Name: fieldStatus,
		Type: fieldTypeSingleSelect,
	},
	{
		Name:      fieldRepo,
		Type:      fieldTypeSingleSelect,
		Provision: true,
		Options: []workflowFieldOption{
			{Name: "vibe-coding-workspace", Description: "Meta workspace repository", Color: "BLUE"},
			{Name: "ww", Description: "Workspace worktree CLI", Color: "GREEN"},
			{Name: "ai-arena", Description: "AI Arena child project", Color: "PURPLE"},
			{Name: "reversi-adventure", Description: "Reversi Adventure child project", Color: "ORANGE"},
			{Name: "vim-learning-game", Description: "Vim Learning Game child project", Color: "PINK"},
			{Name: "envdiff", Description: "envdiff child project", Color: "YELLOW"},
		},
	},
	{
		Name:      fieldKind,
		Type:      fieldTypeSingleSelect,
		Provision: true,
		Options: []workflowFieldOption{
			{Name: "Feature", Description: "Feature work", Color: "GREEN"},
			{Name: "Bug", Description: "Bug fix", Color: "RED"},
			{Name: "Chore", Description: "Maintenance or workflow work", Color: "GRAY"},
			{Name: "Research", Description: "Investigation or discovery work", Color: "BLUE"},
		},
	},
	{
		Name:      fieldPriority,
		Type:      fieldTypeSingleSelect,
		Provision: true,
		Options: []workflowFieldOption{
			{Name: "High", Description: "Needs prompt attention", Color: "RED"},
			{Name: "Medium", Description: "Normal scheduling priority", Color: "YELLOW"},
			{Name: "Low", Description: "Can wait", Color: "BLUE"},
		},
	},
}

var workflowFieldSchemaByName = makeWorkflowFieldSchemaByName()

func makeWorkflowFieldSchemaByName() map[string]workflowFieldSchema {
	byName := make(map[string]workflowFieldSchema, len(workflowFieldSchemas))
	for _, schema := range workflowFieldSchemas {
		byName[schema.Name] = schema
	}
	return byName
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
