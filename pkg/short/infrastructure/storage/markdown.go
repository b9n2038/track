package storage

import (
	"act/pkg/short/domain/model"
	"act/pkg/short/domain/ports/secondary"
	"context"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	ErrLockTimeout = errors.New("timeout waiting for lock")
	lockTimeout    = 5 * time.Second
)

type markdownStorage struct {
	baseDir string
	// locks   map[string]*fileLock
	mu sync.Mutex // Protects the locks map
}

type fileLock struct {
	mu       sync.Mutex
	refCount int
}

func NewMarkdownStorage(baseDir string) secondary.ListStorage {
	return &markdownStorage{
		baseDir: baseDir,
		// locks:   make(map[string]*fileLock),
	}
}

func (s *markdownStorage) getListPath(name string) string {
	return filepath.Join(s.baseDir, ".short", fmt.Sprintf("%s.md", name))
}

func (s *markdownStorage) Exists(name string) bool {
	_, err := os.Stat(s.getListPath(name))
	return err == nil
}
func (s *markdownStorage) acquireLock(lock *fileLock) error {
	ctx, cancel := context.WithTimeout(context.Background(), lockTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		lock.mu.Lock()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("%w: waited %v", ErrLockTimeout, lockTimeout)
	}
}
func (s *markdownStorage) Load(name string) (*model.ShortList, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	path := s.getListPath(name)

	// Default config if file doesn't exist
	config := model.DefaultConfig()
	shortlist := model.NewShortList(name, config)

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return shortlist, nil
		}
		return nil, fmt.Errorf("failed to read list: %w", err)
	}

	sections := s.parseMdWithSections(string(content))
	if sections.Config.MaxCount > 0 {
		shortlist.Config = sections.Config
	}
	shortlist.Open = sections.Open
	shortlist.Closed = sections.Closed

	return shortlist, nil
}

type listSections struct {
	Config model.Config
	Open   []string
	Closed []string
}

func (s *markdownStorage) parseMdWithSections(content string) listSections {
	var sections listSections
	sections.Config = model.DefaultConfig()

	lines := strings.Split(content, "\n")
	currentSection := ""
	inFrontmatter := false
	var frontmatterLines []string

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Handle frontmatter
		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true
				continue
			} else {
				yaml.Unmarshal([]byte(strings.Join(frontmatterLines, "\n")), &sections.Config)
				inFrontmatter = false
				continue
			}
		}

		if inFrontmatter {
			frontmatterLines = append(frontmatterLines, line)
			continue
		}

		// Handle sections
		if strings.HasPrefix(line, "# ") {
			currentSection = strings.TrimPrefix(line, "# ")
			continue
		}

		switch currentSection {
		case "Open":
			sections.Open = append(sections.Open, line)
		case "Closed":
			sections.Closed = append(sections.Closed, line)
		}
	}

	return sections
}

func (s *markdownStorage) Save(list *model.ShortList) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	content := s.createMdWithSections(list)
	path := s.getListPath(list.Name)

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func (s *markdownStorage) createMdWithSections(list *model.ShortList) string {
	var sb strings.Builder

	// Write config as frontmatter
	configBytes, _ := yaml.Marshal(list.Config)
	sb.WriteString("---\n")
	sb.Write(configBytes)
	sb.WriteString("---\n\n")

	// Write open items
	sb.WriteString("# Open\n")
	for _, item := range list.Open {
		sb.WriteString(item)
		sb.WriteString("\n")
	}
	sb.WriteString("\n")

	// Write closed items
	sb.WriteString("# Closed\n")
	for _, item := range list.Closed {
		sb.WriteString(item)
		sb.WriteString("\n")
	}

	return sb.String()
}
