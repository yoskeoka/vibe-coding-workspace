package pj

import "time"

const (
	defaultCachePath = ".local/pj/cache.json"
	fieldStatus      = "Status"
	fieldRepo        = "Repo"
	fieldKind        = "Kind"
	fieldPriority    = "Priority"
)

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
