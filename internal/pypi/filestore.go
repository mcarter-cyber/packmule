package pypi

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

type FileStore struct {
	// IndexPath is the file SaveIndex writes to, e.g. "data/index.json".
	IndexPath string

	// ProjectsDir is the directory SaveProject writes into. If empty,
	// it defaults to "<dir of IndexPath>/projects".
	ProjectsDir string

	// Indent, if non-empty, is passed to json.MarshalIndent (e.g. "  ").
	// Leave empty for compact single-line JSON.
	Indent string

	mu sync.Mutex // guards concurrent writes from SaveProject
}

// NewFileStore builds a FileStore writing the index to indexPath and
// per-project files (if used) to a "projects" directory next to it.
func NewFileStore(indexPath string) *FileStore {
	return &FileStore{
		IndexPath: indexPath,
		Indent:    "  ",
	}
}

// SaveIndex writes idx to fs.IndexPath as JSON, creating parent
// directories as needed. It writes to a temp file first and renames it
// into place, so a crash or interrupted write can't leave a truncated
// index.json behind.
func (fs *FileStore) SaveIndex(ctx context.Context, idx *Index) error {
	if fs.IndexPath == "" {
		return fmt.Errorf("filestore: IndexPath is empty")
	}

	dir := filepath.Dir(fs.IndexPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("filestore: creating dir %q: %w", dir, err)
	}

	data, err := fs.marshal(idx)
	if err != nil {
		return fmt.Errorf("filestore: marshaling index: %w", err)
	}

	return writeFileAtomic(fs.IndexPath, data)
}

// SaveProject writes detail to "<ProjectsDir>/<name>.json". Safe for
// concurrent use, e.g. from Collector's worker pool.
func (fs *FileStore) SaveProject(ctx context.Context, detail *ProjectDetail) error {
	dir := fs.projectsDir()

	fs.mu.Lock()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fs.mu.Unlock()
		return fmt.Errorf("filestore: creating dir %q: %w", dir, err)
	}
	fs.mu.Unlock()

	data, err := fs.marshal(detail)
	if err != nil {
		return fmt.Errorf("filestore: marshaling project %q: %w", detail.Name, err)
	}

	path := filepath.Join(dir, detail.Name+".json")
	return writeFileAtomic(path, data)
}

func (fs *FileStore) projectsDir() string {
	if fs.ProjectsDir != "" {
		return fs.ProjectsDir
	}
	return filepath.Join(filepath.Dir(fs.IndexPath), "projects")
}

func (fs *FileStore) marshal(v interface{}) ([]byte, error) {
	if fs.Indent != "" {
		return json.MarshalIndent(v, "", fs.Indent)
	}
	return json.Marshal(v)
}

// writeFileAtomic writes data to path via a temp file + rename, avoiding
// partial/corrupt files if the process dies mid-write.
func writeFileAtomic(path string, data []byte) error {
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("writing temp file %q: %w", tmp, err)
	}
	if err := os.Rename(tmp, path); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("renaming %q to %q: %w", tmp, path, err)
	}
	return nil
}
