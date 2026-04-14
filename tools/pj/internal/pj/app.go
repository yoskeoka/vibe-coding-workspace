package pj

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

type projectClient interface {
	ensureProject(owner, ownerType, title string) (ProjectRef, bool, error)
	syncProject(ref ProjectRef) (*Cache, error)
	provisionWorkflowFields(cache *Cache) (bool, error)
	addDraftItem(cache *Cache, title, body string, fieldValues map[string]string) error
	moveItem(cache *Cache, itemID, status string) error
}

var newProjectClient = func() (projectClient, error) {
	return newGitHubClient()
}

func Run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		printUsage(stderr)
		return errors.New("missing command")
	}

	switch args[0] {
	case "init":
		return runInit(args[1:], stdout)
	case "sync":
		return runSync(args[1:], stdout)
	case "list":
		return runList(args[1:], stdout)
	case "add":
		return runAdd(args[1:], stdout)
	case "move":
		return runMove(args[1:], stdout)
	case "-h", "--help", "help":
		printUsage(stdout)
		return nil
	default:
		printUsage(stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func runInit(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	owner := fs.String("owner", "", "GitHub login or organization")
	ownerType := fs.String("owner-type", "", "owner type: user or org")
	if err := fs.Parse(args); err != nil {
		return err
	}

	ref := ProjectRef{
		Owner:     *owner,
		OwnerType: *ownerType,
	}
	if cached, err := loadCache(*cachePath); err == nil {
		ref = mergeProjectRef(ref, cached.Project)
	}
	if ref.Owner == "" || ref.OwnerType == "" {
		return fmt.Errorf("init requires --owner and --owner-type, or a cache with owner metadata")
	}

	client, err := newProjectClient()
	if err != nil {
		return err
	}

	projectRef, created, err := client.ensureProject(ref.Owner, ref.OwnerType, canonicalProjectTitle)
	if err != nil {
		return err
	}

	cache, err := client.syncProject(projectRef)
	if err != nil {
		return err
	}
	provisioned, provisionErr := client.provisionWorkflowFields(cache)
	if provisionErr != nil || provisioned {
		refreshed, refreshErr := client.syncProject(cache.Project)
		if refreshErr == nil {
			cache = refreshed
		} else if provisionErr == nil {
			return refreshErr
		}
	}
	if err := writeCache(*cachePath, cache); err != nil {
		return err
	}
	if provisionErr != nil {
		action := "resolved"
		if created {
			action = "created"
		}
		return fmt.Errorf("%s canonical project %q (#%d) and wrote cache to %s, but %w",
			action, cache.Project.Title, cache.Project.ProjectNumber, *cachePath, provisionErr)
	}
	if err := validateRequiredFields(cache); err != nil {
		action := "resolved"
		if created {
			action = "created"
		}
		return fmt.Errorf("%s canonical project %q (#%d) and wrote cache to %s, but %w",
			action, cache.Project.Title, cache.Project.ProjectNumber, *cachePath, err)
	}

	action := "resolved"
	if created {
		action = "created"
	}
	fmt.Fprintf(stdout, "%s canonical project %q (#%d) for %s %q and wrote cache to %s\n",
		action, cache.Project.Title, cache.Project.ProjectNumber, cache.Project.OwnerType, cache.Project.Owner, *cachePath)
	return nil
}

func runSync(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	owner := fs.String("owner", "", "GitHub login or organization")
	ownerType := fs.String("owner-type", "", "owner type: user or org")
	projectNumber := fs.Int("project", 0, "GitHub ProjectV2 number")
	if err := fs.Parse(args); err != nil {
		return err
	}

	ref := ProjectRef{
		Owner:         *owner,
		OwnerType:     *ownerType,
		ProjectNumber: *projectNumber,
	}
	if cached, err := loadCache(*cachePath); err == nil {
		ref = mergeProjectRef(ref, cached.Project)
	}
	if ref.Owner == "" || ref.OwnerType == "" || ref.ProjectNumber == 0 {
		return fmt.Errorf("sync requires --owner, --owner-type, and --project, or a cache with project metadata")
	}

	client, err := newProjectClient()
	if err != nil {
		return err
	}
	cache, err := client.syncProject(ref)
	if err != nil {
		return err
	}
	if err := writeCache(*cachePath, cache); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "synced %d items from %s project %q (#%d) to %s\n",
		len(cache.Items), cache.Project.OwnerType, cache.Project.Owner, cache.Project.ProjectNumber, *cachePath)
	return nil
}

func runList(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("list", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	status := fs.String("status", "", "filter by Status")
	repo := fs.String("repo", "", "filter by Repo")
	kind := fs.String("kind", "", "filter by Kind")
	priority := fs.String("priority", "", "filter by Priority")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cache, err := loadCacheRequired(*cachePath)
	if err != nil {
		return err
	}

	tw := tabwriter.NewWriter(stdout, 0, 4, 2, ' ', 0)
	fmt.Fprintln(tw, "ITEM ID\tSTATUS\tREPO\tKIND\tPRIORITY\tTITLE")
	for _, item := range cache.Items {
		if *status != "" && item.Status != *status {
			continue
		}
		if *repo != "" && item.Repo != *repo && item.Repository != *repo {
			continue
		}
		if *kind != "" && item.Kind != *kind {
			continue
		}
		if *priority != "" && item.Priority != *priority {
			continue
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\t%s\n",
			item.ID, valueOrDash(item.Status), valueOrDash(firstNonEmpty(item.Repo, item.Repository)),
			valueOrDash(item.Kind), valueOrDash(item.Priority), item.Title)
	}
	return tw.Flush()
}

func runAdd(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	title := fs.String("title", "", "draft item title")
	body := fs.String("body", "", "draft item body")
	status := fs.String("status", "", "Status field option")
	repo := fs.String("repo", "", "Repo field option")
	kind := fs.String("kind", "", "Kind field option")
	priority := fs.String("priority", "", "Priority field option")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *title == "" {
		return fmt.Errorf("add requires --title")
	}

	cache, err := loadCacheRequired(*cachePath)
	if err != nil {
		return err
	}
	client, err := newProjectClient()
	if err != nil {
		return err
	}

	fieldValues := map[string]string{
		fieldStatus:   *status,
		fieldRepo:     *repo,
		fieldKind:     *kind,
		fieldPriority: *priority,
	}
	if err := client.addDraftItem(cache, *title, *body, fieldValues); err != nil {
		return err
	}

	refreshed, err := client.syncProject(cache.Project)
	if err != nil {
		return err
	}
	if err := writeCache(*cachePath, refreshed); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "added draft item %q and refreshed cache\n", *title)
	return nil
}

func runMove(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("move", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	itemID := fs.String("item", "", "project item id")
	status := fs.String("status", "", "target Status option")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *itemID == "" || *status == "" {
		return fmt.Errorf("move requires --item and --status")
	}

	cache, err := loadCacheRequired(*cachePath)
	if err != nil {
		return err
	}
	client, err := newProjectClient()
	if err != nil {
		return err
	}
	if err := client.moveItem(cache, *itemID, *status); err != nil {
		return err
	}

	refreshed, err := client.syncProject(cache.Project)
	if err != nil {
		return err
	}
	if err := writeCache(*cachePath, refreshed); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "moved item %s to status %q and refreshed cache\n", *itemID, *status)
	return nil
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, `pj manages a workspace GitHub Project cache.

Usage:
  go -C tools/pj run ./cmd/pj init --owner <owner> --owner-type user|org
  go -C tools/pj run ./cmd/pj sync --owner <owner> --owner-type user|org --project <number>
  go -C tools/pj run ./cmd/pj list [--status <value>] [--repo <value>] [--kind <value>] [--priority <value>]
  go -C tools/pj run ./cmd/pj add --title <title> [--body <text>] [--status <value>] [--repo <value>] [--kind <value>] [--priority <value>]
  go -C tools/pj run ./cmd/pj move --item <item-id> --status <value>

`)
}

func valueOrDash(v string) string {
	if v == "" {
		return "-"
	}
	return v
}

func mergeProjectRef(ref, cached ProjectRef) ProjectRef {
	if ref.Owner == "" {
		ref.Owner = cached.Owner
	}
	if ref.OwnerType == "" {
		ref.OwnerType = cached.OwnerType
	}
	if ref.ProjectNumber == 0 {
		ref.ProjectNumber = cached.ProjectNumber
	}
	return ref
}

func validateRequiredFields(cache *Cache) error {
	missing := make([]string, 0, len(workflowFieldSchemas))
	incompatible := make([]string, 0)
	for _, schema := range workflowFieldSchemas {
		field, ok := cache.Fields[schema.Name]
		if !ok || field.ID == "" {
			missing = append(missing, schema.Name)
			continue
		}
		if err := validateFieldCompatibility(schema, field); err != nil {
			incompatible = append(incompatible, err.Error())
		}
	}
	parts := make([]string, 0, 2)
	if len(missing) > 0 {
		parts = append(parts, fmt.Sprintf("missing required fields: %s", commaList(missing)))
	}
	if len(incompatible) > 0 {
		parts = append(parts, fmt.Sprintf("incompatible required fields: %s", strings.Join(incompatible, "; ")))
	}
	if len(parts) > 0 {
		return fmt.Errorf("project is incompatible with required workflow fields: %s", strings.Join(parts, "; "))
	}
	return nil
}

func validateFieldCompatibility(schema workflowFieldSchema, field FieldCache) error {
	if field.Type != schema.Type {
		return fmt.Errorf("%q has unsupported type %q (want %q)", schema.Name, field.Type, schema.Type)
	}
	if len(schema.Options) == 0 {
		return nil
	}

	missingOptions := make([]string, 0, len(schema.Options))
	for _, option := range schema.Options {
		if field.Options[option.Name] == "" {
			missingOptions = append(missingOptions, option.Name)
		}
	}
	if len(missingOptions) > 0 {
		return fmt.Errorf("%q is missing required options: %s", schema.Name, commaList(missingOptions))
	}
	return nil
}
