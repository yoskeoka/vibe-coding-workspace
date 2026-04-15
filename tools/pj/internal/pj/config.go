package pj

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func loadOwnerConfig(path string) (*OwnerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg OwnerConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("decode config: %w", err)
	}
	if err := validateOwnerConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid owner config at %s: %w; run `pj config set` or `pj config clear` to repair it", path, err)
	}
	return &cfg, nil
}

func writeOwnerConfig(path string, cfg *OwnerConfig) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encode config: %w", err)
	}
	data = append(data, '\n')

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

func writeOwnerConfigFromProject(path string, ref ProjectRef) error {
	return writeOwnerConfig(path, &OwnerConfig{
		Owner:     ref.Owner,
		OwnerType: ref.OwnerType,
	})
}

func resolveOwnerConfig(configPath, cachePath, owner, ownerType string) (*OwnerConfig, error) {
	flagConfig, err := ownerConfigFromFlags(owner, ownerType)
	if err != nil {
		return nil, err
	}

	cfg, err := loadOwnerConfig(configPath)
	if err == nil {
		if flagConfig != nil && !sameOwnerTarget(*cfg, *flagConfig) {
			return nil, fmt.Errorf("owner config mismatch: local config targets %s %q; use `pj config set` or `pj config clear` before switching owner scope", cfg.OwnerType, cfg.Owner)
		}
		return cfg, nil
	}
	if !os.IsNotExist(err) {
		return nil, err
	}

	if flagConfig != nil {
		return flagConfig, nil
	}

	cache, cacheErr := loadCache(cachePath)
	if cacheErr == nil && cache.Project.Owner != "" && cache.Project.OwnerType != "" {
		return &OwnerConfig{
			Owner:     cache.Project.Owner,
			OwnerType: cache.Project.OwnerType,
		}, nil
	}
	if cacheErr != nil && !os.IsNotExist(cacheErr) {
		return nil, cacheErr
	}

	return nil, fmt.Errorf("owner config not found at %s; run `pj init --owner <owner> --owner-type user|org` or `pj config set --owner <owner> --owner-type user|org`", configPath)
}

func ownerConfigFromFlags(owner, ownerType string) (*OwnerConfig, error) {
	if owner == "" && ownerType == "" {
		return nil, nil
	}
	if owner == "" || ownerType == "" {
		return nil, fmt.Errorf("owner flags require both --owner and --owner-type")
	}
	cfg := OwnerConfig{Owner: owner, OwnerType: ownerType}
	if err := validateOwnerConfig(cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func sameOwnerTarget(a, b OwnerConfig) bool {
	return a.Owner == b.Owner && a.OwnerType == b.OwnerType
}

func validateOwnerConfig(cfg OwnerConfig) error {
	if strings.TrimSpace(cfg.Owner) == "" {
		return fmt.Errorf("missing owner")
	}
	switch cfg.OwnerType {
	case "user", "org":
		return nil
	case "":
		return fmt.Errorf("missing owner_type")
	default:
		return fmt.Errorf("unsupported owner_type %q: use user or org", cfg.OwnerType)
	}
}

func removeFileIfExists(path string) error {
	err := os.Remove(path)
	if err == nil || os.IsNotExist(err) {
		return nil
	}
	return err
}
