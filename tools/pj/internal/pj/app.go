package pj

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
)

type projectClient interface {
	ensureProject(owner, ownerType, title string) (ProjectRef, bool, error)
	syncProject(ref ProjectRef) (*Cache, error)
	provisionWorkflowFields(cache *Cache) (bool, error)
	projectLinkedRepositories(projectID string) ([]RepositoryRef, error)
	linkProjectToRepository(projectID string, repo RepositoryRef) error
	unlinkProjectFromRepository(projectID string, repo RepositoryRef) error
	resolveRepository(owner, name string) (RepositoryRef, error)
	addDraftItem(cache *Cache, title, body string, fieldValues map[string]string) error
	updateItem(cache *Cache, itemID string, update itemUpdate) error
}

type app struct {
	newProjectClient func() (projectClient, error)
}

func newApp() app {
	return app{
		newProjectClient: func() (projectClient, error) {
			return newGitHubClient()
		},
	}
}

func (a app) newClient() (projectClient, error) {
	if a.newProjectClient == nil {
		return nil, errors.New("project client factory is not configured")
	}
	return a.newProjectClient()
}

func Run(args []string, stdout, stderr io.Writer) error {
	return newApp().run(args, stdout, stderr)
}

func (a app) run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		printUsage(stderr)
		return errors.New("missing command")
	}

	switch args[0] {
	case "init":
		return a.runInit(args[1:], stdout)
	case "sync":
		return a.runSync(args[1:], stdout)
	case "config":
		return runConfig(args[1:], stdout)
	case "repo-link":
		return a.runRepoLink(args[1:], stdout)
	case "list":
		return runList(args[1:], stdout)
	case "add":
		return a.runAdd(args[1:], stdout)
	case "update":
		return a.runUpdate(args[1:], stdout)
	case "url":
		return runURL(args[1:], stdout)
	case "open":
		return runOpen(args[1:], stdout)
	case "help":
		if len(args) > 1 {
			printCommandUsage(stdout, args[1])
			return nil
		}
		printUsage(stdout)
		return nil
	case "-h", "--help":
		printUsage(stdout)
		return nil
	default:
		printUsage(stderr)
		return fmt.Errorf("unknown command %q", args[0])
	}
}

func (a app) runInit(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("init", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	configPath := fs.String("config", defaultConfigPath, "owner config file path")
	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	owner := fs.String("owner", "", "GitHub login or organization")
	ownerType := fs.String("owner-type", "", "owner type: user or org")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := resolveOwnerConfig(*configPath, *cachePath, *owner, *ownerType)
	if err != nil {
		return err
	}

	client, err := a.newClient()
	if err != nil {
		return err
	}

	projectRef, created, err := client.ensureProject(cfg.Owner, cfg.OwnerType, canonicalProjectTitle)
	if err != nil {
		return err
	}

	cache, err := client.syncProject(projectRef)
	if err != nil {
		return err
	}
	if err := enrichCacheRepoOptions(cache, *cachePath); err != nil {
		return err
	}

	provisioned, provisionErr := client.provisionWorkflowFields(cache)
	var refreshErr error
	if provisionErr != nil || provisioned {
		var refreshed *Cache
		refreshed, refreshErr = client.syncProject(cache.Project)
		if refreshErr == nil {
			cache = refreshed
			if err := enrichCacheRepoOptions(cache, *cachePath); err != nil {
				return err
			}
		} else if provisionErr == nil {
			return refreshErr
		}
	}
	if provisionErr != nil && provisioned && refreshErr != nil {
		action := "resolved"
		if created {
			action = "created"
		}
		return fmt.Errorf("%s canonical project %q (#%d), but provisioning changed remote schema and cache refresh failed; refusing to write stale cache to %s: provision error: %w; refresh error: %v",
			action, cache.Project.Title, cache.Project.ProjectNumber, *cachePath, provisionErr, refreshErr)
	}
	if err := writeOwnerConfigFromProject(*configPath, cache.Project); err != nil {
		return err
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

func (a app) runSync(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("sync", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	configPath := fs.String("config", defaultConfigPath, "owner config file path")
	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	owner := fs.String("owner", "", "GitHub login or organization")
	ownerType := fs.String("owner-type", "", "owner type: user or org")
	projectNumber := fs.Int("project", 0, "GitHub ProjectV2 number")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := resolveOwnerConfig(*configPath, *cachePath, *owner, *ownerType)
	if err != nil {
		return err
	}

	ref := ProjectRef{
		Owner:     cfg.Owner,
		OwnerType: cfg.OwnerType,
	}
	if *projectNumber != 0 {
		ref.ProjectNumber = *projectNumber
	} else if cached, err := loadCache(*cachePath); err == nil {
		ref.ProjectNumber = cached.Project.ProjectNumber
	}
	if ref.ProjectNumber == 0 {
		return fmt.Errorf("sync requires --project, or a cache with project metadata")
	}

	client, err := a.newClient()
	if err != nil {
		return err
	}
	cache, err := client.syncProject(ref)
	if err != nil {
		return err
	}
	if err := enrichCacheRepoOptions(cache, *cachePath); err != nil {
		return err
	}
	if err := writeOwnerConfigFromProject(*configPath, cache.Project); err != nil {
		return err
	}
	if err := writeCache(*cachePath, cache); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "synced %d items from %s project %q (#%d) to %s\n",
		len(cache.Items), cache.Project.OwnerType, cache.Project.Owner, cache.Project.ProjectNumber, *cachePath)
	return nil
}

func runConfig(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("config requires a subcommand: show, set, or clear")
	}

	switch args[0] {
	case "show":
		return runConfigShow(args[1:], stdout)
	case "set":
		return runConfigSet(args[1:], stdout)
	case "clear":
		return runConfigClear(args[1:], stdout)
	default:
		return fmt.Errorf("unknown config subcommand %q", args[0])
	}
}

func runConfigShow(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("config show", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	configPath := fs.String("config", defaultConfigPath, "owner config file path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := loadOwnerConfig(*configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("owner config not found at %s; run `pj init` or `pj config set` first", *configPath)
		}
		return err
	}

	fmt.Fprintf(stdout, "%s %q\n", cfg.OwnerType, cfg.Owner)
	return nil
}

func runConfigSet(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("config set", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	configPath := fs.String("config", defaultConfigPath, "owner config file path")
	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	owner := fs.String("owner", "", "GitHub login or organization")
	ownerType := fs.String("owner-type", "", "owner type: user or org")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := ownerConfigFromFlags(*owner, *ownerType)
	if err != nil {
		return err
	}
	if cfg == nil {
		return fmt.Errorf("config set requires --owner and --owner-type")
	}

	if cached, err := loadCache(*cachePath); err == nil {
		cachedOwner := OwnerConfig{Owner: cached.Project.Owner, OwnerType: cached.Project.OwnerType}
		if cachedOwner.Owner != "" && cachedOwner.OwnerType != "" && !sameOwnerTarget(*cfg, cachedOwner) {
			if err := removeFileIfExists(*cachePath); err != nil {
				return fmt.Errorf("clear stale cache: %w", err)
			}
		}
	} else if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := writeOwnerConfig(*configPath, cfg); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "set owner config to %s %q\n", cfg.OwnerType, cfg.Owner)
	return nil
}

func runConfigClear(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("config clear", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	configPath := fs.String("config", defaultConfigPath, "owner config file path")
	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := removeFileIfExists(*configPath); err != nil {
		return fmt.Errorf("clear owner config: %w", err)
	}
	if err := removeFileIfExists(*cachePath); err != nil {
		return fmt.Errorf("clear cache: %w", err)
	}
	fmt.Fprintln(stdout, "cleared owner config and cache")
	return nil
}

func (a app) runRepoLink(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("repo-link requires a subcommand: status, add, or remove")
	}

	switch args[0] {
	case "status":
		return a.runRepoLinkStatus(args[1:], stdout)
	case "add":
		return a.runRepoLinkAdd(args[1:], stdout)
	case "remove":
		return a.runRepoLinkRemove(args[1:], stdout)
	default:
		return fmt.Errorf("unknown repo-link subcommand %q", args[0])
	}
}

func (a app) runRepoLinkStatus(args []string, stdout io.Writer) error {
	cache, repo, client, err := a.repoLinkInputs("repo-link status", args)
	if err != nil {
		return err
	}

	linked, err := projectHasRepositoryLink(client, cache.Project.ProjectID, repo.NameWithOwner)
	if err != nil {
		return err
	}
	if linked {
		fmt.Fprintf(stdout, "%s is linked to project %q (#%d)\n", repo.NameWithOwner, cache.Project.Title, cache.Project.ProjectNumber)
		return nil
	}
	fmt.Fprintf(stdout, "%s is not linked to project %q (#%d)\n", repo.NameWithOwner, cache.Project.Title, cache.Project.ProjectNumber)
	return nil
}

func (a app) runRepoLinkAdd(args []string, stdout io.Writer) error {
	cache, repo, client, err := a.repoLinkInputs("repo-link add", args)
	if err != nil {
		return err
	}

	linked, err := projectHasRepositoryLink(client, cache.Project.ProjectID, repo.NameWithOwner)
	if err != nil {
		return err
	}
	if linked {
		fmt.Fprintf(stdout, "%s is already linked to project %q (#%d)\n", repo.NameWithOwner, cache.Project.Title, cache.Project.ProjectNumber)
		return nil
	}
	if err := client.linkProjectToRepository(cache.Project.ProjectID, repo); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "linked %s to project %q (#%d)\n", repo.NameWithOwner, cache.Project.Title, cache.Project.ProjectNumber)
	return nil
}

func (a app) runRepoLinkRemove(args []string, stdout io.Writer) error {
	cache, repo, client, err := a.repoLinkInputs("repo-link remove", args)
	if err != nil {
		return err
	}

	linked, err := projectHasRepositoryLink(client, cache.Project.ProjectID, repo.NameWithOwner)
	if err != nil {
		return err
	}
	if !linked {
		fmt.Fprintf(stdout, "%s is not linked to project %q (#%d)\n", repo.NameWithOwner, cache.Project.Title, cache.Project.ProjectNumber)
		return nil
	}
	if err := client.unlinkProjectFromRepository(cache.Project.ProjectID, repo); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "unlinked %s from project %q (#%d)\n", repo.NameWithOwner, cache.Project.Title, cache.Project.ProjectNumber)
	return nil
}

func (a app) repoLinkInputs(command string, args []string) (*Cache, RepositoryRef, projectClient, error) {
	fs := flag.NewFlagSet(command, flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	if err := fs.Parse(args); err != nil {
		return nil, RepositoryRef{}, nil, err
	}
	if fs.NArg() != 1 {
		return nil, RepositoryRef{}, nil, fmt.Errorf("%s requires exactly one <owner>/<repo> argument", command)
	}

	target, err := parseRepositoryTarget(fs.Arg(0))
	if err != nil {
		return nil, RepositoryRef{}, nil, err
	}
	cache, err := loadCacheRequired(*cachePath)
	if err != nil {
		return nil, RepositoryRef{}, nil, err
	}
	if cache.Project.ProjectID == "" {
		return nil, RepositoryRef{}, nil, fmt.Errorf("cache is missing project id; run `pj init` or `pj sync`")
	}
	if cache.Project.Owner == "" {
		return nil, RepositoryRef{}, nil, fmt.Errorf("cache is missing project owner; run `pj init` or `pj sync`")
	}
	if target.Owner != cache.Project.Owner {
		return nil, RepositoryRef{}, nil, fmt.Errorf("repository owner %q does not match project owner %q; ProjectV2 repository links require the same owner", target.Owner, cache.Project.Owner)
	}

	client, err := a.newClient()
	if err != nil {
		return nil, RepositoryRef{}, nil, err
	}
	repo, err := client.resolveRepository(target.Owner, target.Name)
	if err != nil {
		return nil, RepositoryRef{}, nil, err
	}
	return cache, repo, client, nil
}

func parseRepositoryTarget(value string) (repositoryTarget, error) {
	parts := strings.Split(strings.TrimSpace(value), "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return repositoryTarget{}, fmt.Errorf("repository target must be <owner>/<repo>")
	}
	return repositoryTarget{Owner: parts[0], Name: parts[1]}, nil
}

func projectHasRepositoryLink(client projectClient, projectID, nameWithOwner string) (bool, error) {
	repos, err := client.projectLinkedRepositories(projectID)
	if err != nil {
		return false, err
	}
	for _, linkedRepo := range repos {
		if linkedRepo.NameWithOwner == nameWithOwner {
			return true, nil
		}
	}
	return false, nil
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

func (a app) runAdd(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("add", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	title := fs.String("title", "", "draft item title")
	body := fs.String("body", "", "draft item body")
	bodyFile := fs.String("body-file", "", "path to a file containing the draft item body")
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

	resolvedBody, err := resolveBodyInput("add", *body, *bodyFile, flagWasProvided(fs, "body"), flagWasProvided(fs, "body-file"))
	if err != nil {
		return err
	}
	cache, err := loadCacheRequired(*cachePath)
	if err != nil {
		return err
	}

	fieldValues, err := resolveFieldInputs(cache, map[string]string{
		fieldStatus:   *status,
		fieldRepo:     *repo,
		fieldKind:     *kind,
		fieldPriority: *priority,
	})
	if err != nil {
		return err
	}
	client, err := a.newClient()
	if err != nil {
		return err
	}
	if err := client.addDraftItem(cache, *title, resolvedBody, fieldValues); err != nil {
		return err
	}

	refreshed, err := client.syncProject(cache.Project)
	if err != nil {
		return err
	}
	if err := enrichCacheRepoOptions(refreshed, *cachePath); err != nil {
		return err
	}
	if err := writeCache(*cachePath, refreshed); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "added draft item %q and refreshed cache\n", *title)
	return nil
}

func (a app) runUpdate(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("update", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	itemID := fs.String("item", "", "project item id")
	title := fs.String("title", "", "draft item title")
	body := fs.String("body", "", "draft item body")
	bodyFile := fs.String("body-file", "", "path to a file containing the draft item body")
	status := fs.String("status", "", "Status field option")
	repo := fs.String("repo", "", "Repo field option")
	kind := fs.String("kind", "", "Kind field option")
	priority := fs.String("priority", "", "Priority field option")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *itemID == "" {
		return fmt.Errorf("update requires --item")
	}

	inlineBodyProvided := flagWasProvided(fs, "body")
	bodyFileProvided := flagWasProvided(fs, "body-file")
	bodyProvided := inlineBodyProvided || bodyFileProvided
	resolvedBody, err := resolveBodyInput("update", *body, *bodyFile, inlineBodyProvided, bodyFileProvided)
	if err != nil {
		return err
	}

	cache, err := loadCacheRequired(*cachePath)
	if err != nil {
		return err
	}

	fieldValues, err := resolveFieldInputs(cache, map[string]string{
		fieldStatus:   *status,
		fieldRepo:     *repo,
		fieldKind:     *kind,
		fieldPriority: *priority,
	})
	if err != nil {
		return err
	}
	update := itemUpdate{
		Title:         *title,
		TitleProvided: flagWasProvided(fs, "title"),
		Body:          resolvedBody,
		BodyProvided:  bodyProvided,
		FieldValues:   fieldValues,
	}
	if !update.TitleProvided && !update.BodyProvided && len(update.FieldValues) == 0 {
		return fmt.Errorf("update requires at least one field: --title, --body, --body-file, --status, --repo, --kind, or --priority")
	}
	client, err := a.newClient()
	if err != nil {
		return err
	}
	if err := client.updateItem(cache, *itemID, update); err != nil {
		return err
	}

	refreshed, err := client.syncProject(cache.Project)
	if err != nil {
		return err
	}
	if err := enrichCacheRepoOptions(refreshed, *cachePath); err != nil {
		return err
	}
	if err := writeCache(*cachePath, refreshed); err != nil {
		return err
	}

	fmt.Fprintf(stdout, "updated item %s and refreshed cache\n", *itemID)
	return nil
}

func runURL(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("url", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cache, err := loadCacheRequired(*cachePath)
	if err != nil {
		return err
	}
	url, err := projectURL(cache.Project)
	if err != nil {
		return err
	}
	fmt.Fprintln(stdout, url)
	return nil
}

var openProjectURL = openURLInBrowser

func runOpen(args []string, stdout io.Writer) error {
	fs := flag.NewFlagSet("open", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	cachePath := fs.String("cache", defaultCachePath, "cache file path")
	if err := fs.Parse(args); err != nil {
		return err
	}

	cache, err := loadCacheRequired(*cachePath)
	if err != nil {
		return err
	}
	url, err := projectURL(cache.Project)
	if err != nil {
		return err
	}
	if err := openProjectURL(url); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "opened %s\n", url)
	return nil
}

func printUsage(w io.Writer) {
	fmt.Fprintf(w, `pj manages a workspace GitHub Project cache.

Usage:
  go -C tools/pj run ./cmd/pj init --owner <owner> --owner-type user|org
  go -C tools/pj run ./cmd/pj sync [--project <number>]
  go -C tools/pj run ./cmd/pj config show
  go -C tools/pj run ./cmd/pj config set --owner <owner> --owner-type user|org
  go -C tools/pj run ./cmd/pj config clear
  go -C tools/pj run ./cmd/pj repo-link status <owner>/<repo>
  go -C tools/pj run ./cmd/pj repo-link add <owner>/<repo>
  go -C tools/pj run ./cmd/pj repo-link remove <owner>/<repo>
  go -C tools/pj run ./cmd/pj list [--status <value>] [--repo <value>] [--kind <value>] [--priority <value>]
  go -C tools/pj run ./cmd/pj add --title <title> [--body <text>|--body-file <path>] [--status <value>] [--repo <value>] [--kind <value>] [--priority <value>]
  go -C tools/pj run ./cmd/pj update --item <item-id> [--title <title>] [--body <text>|--body-file <path>] [--status <value>] [--repo <value>] [--kind <value>] [--priority <value>]
  go -C tools/pj run ./cmd/pj url
  go -C tools/pj run ./cmd/pj open

`)
}

func printCommandUsage(w io.Writer, command string) {
	switch command {
	case "add":
		fmt.Fprintln(w, "Usage: pj add --title <title> [--body <text>|--body-file <path>] [--status <value>] [--repo <value>] [--kind <value>] [--priority <value>]")
	case "update":
		fmt.Fprintln(w, "Usage: pj update --item <item-id> [--title <title>] [--body <text>|--body-file <path>] [--status <value>] [--repo <value>] [--kind <value>] [--priority <value>]")
	case "url":
		fmt.Fprintln(w, "Usage: pj url")
	case "open":
		fmt.Fprintln(w, "Usage: pj open")
	case "init":
		fmt.Fprintln(w, "Usage: pj init --owner <owner> --owner-type user|org")
	case "sync":
		fmt.Fprintln(w, "Usage: pj sync [--project <number>]")
	case "list":
		fmt.Fprintln(w, "Usage: pj list [--status <value>] [--repo <value>] [--kind <value>] [--priority <value>]")
	case "config":
		fmt.Fprintln(w, "Usage: pj config show|set|clear")
	case "repo-link":
		fmt.Fprintln(w, "Usage: pj repo-link status|add|remove <owner>/<repo>")
	default:
		printUsage(w)
	}
}

type itemUpdate struct {
	Title         string
	TitleProvided bool
	Body          string
	BodyProvided  bool
	FieldValues   map[string]string
}

func resolveBodyInput(command, body, bodyFile string, bodyProvided, bodyFileProvided bool) (string, error) {
	if bodyProvided && bodyFileProvided {
		return "", fmt.Errorf("%s accepts either --body or --body-file, not both", command)
	}
	if !bodyFileProvided {
		return body, nil
	}
	data, err := os.ReadFile(bodyFile)
	if err != nil {
		return "", fmt.Errorf("read --body-file %s: %w", bodyFile, err)
	}
	return string(data), nil
}

func flagWasProvided(fs *flag.FlagSet, name string) bool {
	provided := false
	fs.Visit(func(f *flag.Flag) {
		if f.Name == name {
			provided = true
		}
	})
	return provided
}

func resolveFieldInputs(cache *Cache, inputs map[string]string) (map[string]string, error) {
	resolved := map[string]string{}
	for fieldName, value := range inputs {
		if value == "" {
			continue
		}
		var out string
		var err error
		if fieldName == fieldRepo {
			out, err = resolveRepoValue(cache, value)
		} else {
			out, err = resolveSingleSelectValue(cache, fieldName, value)
		}
		if err != nil {
			return nil, err
		}
		resolved[fieldName] = out
	}
	return resolved, nil
}

func resolveSingleSelectValue(cache *Cache, fieldName, input string) (string, error) {
	field, ok := cache.Fields[fieldName]
	if !ok || field.ID == "" {
		return "", fmt.Errorf("project field %q is missing from cache; run `pj sync` against a board with that field", fieldName)
	}
	candidates := make([]valueCandidate, 0, len(field.Options))
	for display := range field.Options {
		candidates = append(candidates, valueCandidate{Display: display, Keys: []string{display}})
	}
	return resolveValue(fieldName, input, candidates)
}

func resolveRepoValue(cache *Cache, input string) (string, error) {
	if len(cache.RepoOptions) == 0 {
		return "", fmt.Errorf("repo metadata is missing from cache; run `pj init` or `pj sync`")
	}
	if index, err := strconv.Atoi(input); err == nil {
		if index < 1 || index > len(cache.RepoOptions) {
			return "", fmt.Errorf("repo index %d is out of range; valid range is 1-%d", index, len(cache.RepoOptions))
		}
		return cache.RepoOptions[index-1].DisplayValue, nil
	}

	candidates := make([]valueCandidate, 0, len(cache.RepoOptions))
	for _, option := range cache.RepoOptions {
		keys := []string{option.CanonicalSlug, option.DisplayValue}
		keys = append(keys, option.Aliases...)
		candidates = append(candidates, valueCandidate{
			Display: option.DisplayValue,
			Keys:    keys,
		})
	}
	return resolveValue(fieldRepo, input, candidates)
}

type valueCandidate struct {
	Display string
	Keys    []string
}

func resolveValue(fieldName, input string, candidates []valueCandidate) (string, error) {
	normalizedInput := normalizeValue(input)
	exact := make(map[string]string)
	for _, candidate := range candidates {
		for _, key := range candidate.Keys {
			if normalizeValue(key) == normalizedInput {
				exact[candidate.Display] = candidate.Display
			}
		}
	}
	if len(exact) == 1 {
		for display := range exact {
			return display, nil
		}
	}
	if len(exact) > 1 {
		return "", fmt.Errorf("ambiguous %s value %q matches: %s", fieldName, input, sortedMapKeys(exact))
	}

	prefix := make(map[string]string)
	for _, candidate := range candidates {
		for _, key := range candidate.Keys {
			if strings.HasPrefix(normalizeValue(key), normalizedInput) {
				prefix[candidate.Display] = candidate.Display
			}
		}
	}
	if len(prefix) == 1 {
		for display := range prefix {
			return display, nil
		}
	}
	if len(prefix) > 1 {
		return "", fmt.Errorf("ambiguous %s value %q matches: %s", fieldName, input, sortedMapKeys(prefix))
	}
	return "", fmt.Errorf("unknown %s value %q", fieldName, input)
}

func normalizeValue(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	var b strings.Builder
	lastSep := false
	for _, r := range value {
		switch r {
		case '-', '_', ' ', '\t', '\n', '\r':
			if !lastSep {
				b.WriteByte(' ')
				lastSep = true
			}
		default:
			b.WriteRune(r)
			lastSep = false
		}
	}
	return strings.TrimSpace(b.String())
}

func sortedMapKeys(values map[string]string) string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return strings.Join(keys, ", ")
}

func projectURL(ref ProjectRef) (string, error) {
	if ref.Owner == "" || ref.ProjectNumber == 0 {
		return "", fmt.Errorf("cache is missing project owner or number; run `pj init` or `pj sync`")
	}
	switch ref.OwnerType {
	case "user":
		return fmt.Sprintf("https://github.com/users/%s/projects/%d", ref.Owner, ref.ProjectNumber), nil
	case "org":
		return fmt.Sprintf("https://github.com/orgs/%s/projects/%d", ref.Owner, ref.ProjectNumber), nil
	default:
		return "", fmt.Errorf("unsupported owner type %q: use user or org", ref.OwnerType)
	}
}

func openURLInBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("open browser for %s: %w", url, err)
	}
	if err := cmd.Process.Release(); err != nil {
		return fmt.Errorf("release browser opener process for %s: %w", url, err)
	}
	return nil
}

func valueOrDash(v string) string {
	if v == "" {
		return "-"
	}
	return v
}

func validateRequiredFields(cache *Cache) error {
	schemas := workflowFieldSchemasForCache(cache)
	missing := make([]string, 0, len(schemas))
	incompatible := make([]string, 0)
	for _, schema := range schemas {
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
