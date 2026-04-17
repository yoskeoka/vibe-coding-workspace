package pj

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

const repoSourceGitHub = "github-repo"

func enrichCacheRepoOptions(cache *Cache, cachePath string) error {
	options, err := deriveWorkspaceRepoOptions(cache, cachePath)
	if err != nil {
		return err
	}
	cache.RepoOptions = options
	return nil
}

func deriveWorkspaceRepoOptions(cache *Cache, cachePath string) ([]RepoOption, error) {
	root := workspaceRootForLocalPath(cachePath)
	rawURLs, err := readSetupRepoURLs(filepath.Join(root, "setup.sh"))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		return fallbackRepoOptions(cache), nil
	}

	if cache.Project.Owner != "" {
		rawURLs = append(rawURLs, fmt.Sprintf("https://github.com/%s/vibe-coding-workspace", cache.Project.Owner))
	}

	seen := map[string]bool{}
	options := make([]RepoOption, 0, len(rawURLs))
	for _, rawURL := range rawURLs {
		option, err := repoOptionFromURL(rawURL)
		if err != nil {
			return nil, err
		}
		if seen[option.CanonicalSlug] {
			continue
		}
		seen[option.CanonicalSlug] = true
		options = append(options, option)
	}

	sort.Slice(options, func(i, j int) bool {
		return options[i].CanonicalSlug < options[j].CanonicalSlug
	})
	return options, nil
}

func fallbackRepoOptions(cache *Cache) []RepoOption {
	owner := firstNonEmpty(cache.Project.Owner, "unknown")
	displays := make([]string, 0)
	if field, ok := cache.Fields[fieldRepo]; ok {
		for display := range field.Options {
			displays = append(displays, display)
		}
	}
	if len(displays) == 0 {
		for _, option := range workflowFieldSchemaByName[fieldRepo].Options {
			displays = append(displays, option.Name)
		}
	}
	sort.Strings(displays)
	options := make([]RepoOption, 0, len(displays))
	for _, display := range displays {
		options = append(options, RepoOption{
			DisplayValue:  display,
			SourceType:    repoSourceGitHub,
			SourceURL:     fmt.Sprintf("https://github.com/%s/%s", owner, display),
			CanonicalSlug: fmt.Sprintf("github.com/%s/%s", owner, display),
			Aliases:       []string{display, fmt.Sprintf("%s/%s", owner, display)},
		})
	}
	return options
}

func workspaceRootForLocalPath(localPath string) string {
	if root, ok := findGitRoot("."); ok {
		return root
	}
	if filepath.IsAbs(localPath) {
		return filepath.Dir(filepath.Dir(filepath.Dir(localPath)))
	}
	return "."
}

func readSetupRepoURLs(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read setup repo metadata from %s: %w", path, err)
	}

	re := regexp.MustCompile(`(?m)^\s*"([^"]+)"\s*$`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	urls := make([]string, 0, len(matches))
	for _, match := range matches {
		urls = append(urls, match[1])
	}
	return urls, nil
}

func repoOptionFromURL(rawURL string) (RepoOption, error) {
	owner, repo, err := githubOwnerRepo(rawURL)
	if err != nil {
		return RepoOption{}, err
	}
	display := strings.TrimSuffix(repo, ".git")
	sourceURL := fmt.Sprintf("https://github.com/%s/%s", owner, display)
	return RepoOption{
		DisplayValue:  display,
		SourceType:    repoSourceGitHub,
		SourceURL:     sourceURL,
		CanonicalSlug: fmt.Sprintf("github.com/%s/%s", owner, display),
		Aliases:       []string{display, fmt.Sprintf("%s/%s", owner, display)},
	}, nil
}

func githubOwnerRepo(rawURL string) (string, string, error) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", "", fmt.Errorf("empty repo URL")
	}

	if strings.HasPrefix(rawURL, "git@github.com:") {
		parts := strings.Split(strings.TrimPrefix(rawURL, "git@github.com:"), "/")
		if len(parts) == 2 {
			return parts[0], strings.TrimSuffix(parts[1], ".git"), nil
		}
	}

	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Host != "github.com" {
		return "", "", fmt.Errorf("unsupported repo URL %q", rawURL)
	}
	parts := strings.Split(strings.Trim(path.Clean(parsed.Path), "/"), "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unsupported repo URL %q", rawURL)
	}
	return parts[0], strings.TrimSuffix(parts[1], ".git"), nil
}
