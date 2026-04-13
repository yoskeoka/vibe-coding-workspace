package pj

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

const githubGraphQLEndpoint = "https://api.github.com/graphql"

type githubClient struct {
	httpClient *http.Client
	token      string
}

func newGitHubClient() (*githubClient, error) {
	out, err := exec.Command("gh", "auth", "token").CombinedOutput()
	if err != nil {
		msg := strings.TrimSpace(string(out))
		if msg == "" {
			msg = err.Error()
		}
		return nil, fmt.Errorf("gh auth token failed: %s", msg)
	}

	token := strings.TrimSpace(string(out))
	if token == "" {
		return nil, fmt.Errorf("gh auth token returned an empty token")
	}

	return &githubClient{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
	}, nil
}

func (c *githubClient) graphQL(query string, variables map[string]any, respData any) error {
	payload := map[string]any{
		"query":     query,
		"variables": variables,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode graphql request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, githubGraphQLEndpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("build graphql request: %w", err)
	}
	req.Header.Set("Authorization", "bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post graphql request: %w", err)
	}
	defer res.Body.Close()

	var gqlResp struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}

	if err := json.NewDecoder(res.Body).Decode(&gqlResp); err != nil {
		return fmt.Errorf("decode graphql response: %w", err)
	}
	if res.StatusCode >= 400 {
		return fmt.Errorf("graphql http status %s", res.Status)
	}
	if len(gqlResp.Errors) > 0 {
		msgs := make([]string, 0, len(gqlResp.Errors))
		for _, err := range gqlResp.Errors {
			msgs = append(msgs, err.Message)
		}
		return fmt.Errorf("graphql error: %s", strings.Join(msgs, "; "))
	}
	if respData == nil {
		return nil
	}
	if err := json.Unmarshal(gqlResp.Data, respData); err != nil {
		return fmt.Errorf("decode graphql data: %w", err)
	}
	return nil
}

func (c *githubClient) syncProject(ref ProjectRef) (*Cache, error) {
	query, rootKey, err := projectQuery(ref.OwnerType)
	if err != nil {
		return nil, err
	}

	var resp map[string]json.RawMessage
	if err := c.graphQL(query, map[string]any{
		"login":  ref.Owner,
		"number": ref.ProjectNumber,
	}, &resp); err != nil {
		return nil, err
	}

	rawRoot, ok := resp[rootKey]
	if !ok {
		return nil, fmt.Errorf("graphql response missing %q root", rootKey)
	}

	var root struct {
		ProjectV2 *projectNode `json:"projectV2"`
	}
	if err := json.Unmarshal(rawRoot, &root); err != nil {
		return nil, fmt.Errorf("decode project root: %w", err)
	}
	if root.ProjectV2 == nil {
		return nil, fmt.Errorf("project not found for %s %q number %d", ref.OwnerType, ref.Owner, ref.ProjectNumber)
	}

	fields := make(map[string]FieldCache)
	for _, node := range root.ProjectV2.Fields.Nodes {
		if node.Name == "" {
			continue
		}
		field := FieldCache{
			ID:      node.ID,
			Type:    node.Typename,
			Options: map[string]string{},
		}
		for _, opt := range node.Options {
			field.Options[opt.Name] = opt.ID
		}
		if len(field.Options) == 0 {
			field.Options = nil
		}
		fields[node.Name] = field
	}

	items := make([]Item, 0, len(root.ProjectV2.Items.Nodes))
	for _, node := range root.ProjectV2.Items.Nodes {
		item := Item{
			ID:          node.ID,
			ContentType: node.Content.Typename,
			Title:       firstNonEmpty(node.Content.Title, "(untitled)"),
			Body:        node.Content.Body,
			URL:         node.Content.URL,
			Repository:  node.Content.Repository.NameWithOwner,
		}
		if item.ContentType == "" {
			item.ContentType = "Unknown"
		}

		for _, fv := range node.FieldValues.Nodes {
			switch fv.Field.Name {
			case fieldStatus:
				item.Status = fv.Name
			case fieldRepo:
				item.Repo = firstNonEmpty(fv.Name, fv.Text)
			case fieldKind:
				item.Kind = firstNonEmpty(fv.Name, fv.Text)
			case fieldPriority:
				item.Priority = firstNonEmpty(fv.Name, fv.Text)
			}
		}

		items = append(items, item)
	}

	cache := &Cache{
		SyncedAt: time.Now().UTC(),
		Project: ProjectRef{
			Owner:         ref.Owner,
			OwnerType:     ref.OwnerType,
			ProjectNumber: ref.ProjectNumber,
			ProjectID:     root.ProjectV2.ID,
			Title:         root.ProjectV2.Title,
		},
		Fields: fields,
		Items:  items,
	}

	return cache, nil
}

func (c *githubClient) addDraftItem(cache *Cache, title, body string, fieldValues map[string]string) error {
	var resp struct {
		AddProjectV2DraftIssue struct {
			ProjectItem struct {
				ID string `json:"id"`
			} `json:"projectItem"`
		} `json:"addProjectV2DraftIssue"`
	}

	err := c.graphQL(`
mutation($projectId: ID!, $title: String!, $body: String!) {
  addProjectV2DraftIssue(input: {projectId: $projectId, title: $title, body: $body}) {
    projectItem {
      id
    }
  }
}
`, map[string]any{
		"projectId": cache.Project.ProjectID,
		"title":     title,
		"body":      body,
	}, &resp)
	if err != nil {
		return err
	}

	itemID := resp.AddProjectV2DraftIssue.ProjectItem.ID
	if itemID == "" {
		return fmt.Errorf("draft item mutation returned an empty item id")
	}

	for fieldName, optionName := range fieldValues {
		if optionName == "" {
			continue
		}
		if err := c.setSingleSelectField(cache, itemID, fieldName, optionName); err != nil {
			return err
		}
	}

	return nil
}

func (c *githubClient) moveItem(cache *Cache, itemID, status string) error {
	return c.setSingleSelectField(cache, itemID, fieldStatus, status)
}

func (c *githubClient) setSingleSelectField(cache *Cache, itemID, fieldName, optionName string) error {
	field, ok := cache.Fields[fieldName]
	if !ok || field.ID == "" {
		return fmt.Errorf("project field %q is missing from cache; run `pj sync` against a board with that field", fieldName)
	}
	optionID := field.Options[optionName]
	if optionID == "" {
		return fmt.Errorf("field %q does not contain option %q", fieldName, optionName)
	}

	return c.graphQL(`
mutation($projectId: ID!, $itemId: ID!, $fieldId: ID!, $optionId: String!) {
  updateProjectV2ItemFieldValue(input: {
    projectId: $projectId
    itemId: $itemId
    fieldId: $fieldId
    value: {singleSelectOptionId: $optionId}
  }) {
    projectV2Item {
      id
    }
  }
}
`, map[string]any{
		"projectId": cache.Project.ProjectID,
		"itemId":    itemID,
		"fieldId":   field.ID,
		"optionId":  optionID,
	}, nil)
}

type projectNode struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Fields struct {
		Nodes []struct {
			Typename string `json:"__typename"`
			ID       string `json:"id"`
			Name     string `json:"name"`
			Options  []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"options"`
		} `json:"nodes"`
	} `json:"fields"`
	Items struct {
		Nodes []struct {
			ID      string `json:"id"`
			Content struct {
				Typename   string `json:"__typename"`
				Title      string `json:"title"`
				Body       string `json:"body"`
				URL        string `json:"url"`
				Repository struct {
					NameWithOwner string `json:"nameWithOwner"`
				} `json:"repository"`
			} `json:"content"`
			FieldValues struct {
				Nodes []struct {
					Typename string `json:"__typename"`
					Name     string `json:"name"`
					Text     string `json:"text"`
					Field    struct {
						Name string `json:"name"`
					} `json:"field"`
				} `json:"nodes"`
			} `json:"fieldValues"`
		} `json:"nodes"`
	} `json:"items"`
}

func projectQuery(ownerType string) (query string, rootKey string, err error) {
	const projectBody = `
projectV2(number: $number) {
  id
  title
  fields(first: 50) {
    nodes {
      __typename
      ... on ProjectV2FieldCommon {
        id
        name
      }
      ... on ProjectV2SingleSelectField {
        options {
          id
          name
        }
      }
    }
  }
  items(first: 100) {
    nodes {
      id
      content {
        __typename
        ... on DraftIssue {
          title
          body
        }
        ... on Issue {
          title
          url
          repository {
            nameWithOwner
          }
        }
        ... on PullRequest {
          title
          url
          repository {
            nameWithOwner
          }
        }
      }
      fieldValues(first: 20) {
        nodes {
          __typename
          ... on ProjectV2ItemFieldSingleSelectValue {
            name
            field {
              ... on ProjectV2FieldCommon {
                name
              }
            }
          }
          ... on ProjectV2ItemFieldTextValue {
            text
            field {
              ... on ProjectV2FieldCommon {
                name
              }
            }
          }
        }
      }
    }
  }
}`

	switch ownerType {
	case "user":
		return fmt.Sprintf(`
query($login: String!, $number: Int!) {
  user(login: $login) {
    %s
  }
}`, projectBody), "user", nil
	case "org":
		return fmt.Sprintf(`
query($login: String!, $number: Int!) {
  organization(login: $login) {
    %s
  }
}`, projectBody), "organization", nil
	default:
		return "", "", fmt.Errorf("unsupported owner type %q: use user or org", ownerType)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
