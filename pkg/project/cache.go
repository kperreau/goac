package project

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

type Cache struct {
	Target map[Target]*Metadata
	Path   string
}

var DefaultCachePath = ".goac/cache/"

func (p *Project) LoadCache() error {
	cacheFilePath := fmt.Sprintf("%s%s.yaml", DefaultCachePath, p.HashPath)

	// init a default basic cache
	cacheData := Cache{Path: p.CleanPath, Target: map[Target]*Metadata{}}

	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		p.Cache = &cacheData
		return nil
	}

	if err := readCacheFromFile(cacheFilePath, &cacheData); err != nil {
		return fmt.Errorf("error loading cache: %w", err)
	}

	p.Cache = &cacheData

	return nil
}

func readCacheFromFile(path string, cache *Cache) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(data, cache); err != nil {
		return fmt.Errorf("error unmarshaling cache data: %w", err)
	}

	return nil
}

func (cm *Metadata) isMetadataMatch(m *Metadata) bool {
	return cm.DependenciesHash == m.DependenciesHash &&
		cm.DirHash == m.DirHash
}

func (p *Project) writeCache() error {
	cacheFilePath := filepath.Join(DefaultCachePath, fmt.Sprintf("%s.yaml", p.HashPath))

	if err := os.MkdirAll(DefaultCachePath, 0o755); err != nil {
		return fmt.Errorf("error creating cache directory: %v", err)
	}

	p.Cache.Target[p.CMDOptions.Target] = &Metadata{
		DependenciesHash: p.Metadata.DependenciesHash,
		DirHash:          p.Metadata.DirHash,
		Date:             time.Now().Format(time.RFC3339),
	}

	cacheData, err := yaml.Marshal(p.Cache)
	if err != nil {
		return fmt.Errorf("error encoding yaml data: %v", err)
	}

	if err := os.WriteFile(cacheFilePath, cacheData, 0o644); err != nil {
		return fmt.Errorf("error writing cache file %s: %v", cacheFilePath, err)
	}
	return nil
}
