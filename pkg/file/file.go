package file

import (
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	Path string
}

func New(path string) *File {
	return &File{
		Path: path,
	}
}

func (f *File) Read() ([]byte, error) {
	data, err := os.ReadFile(f.Path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// WriteFile writes data to a file inside dir, creating the directory if needed.
func WriteFile(dir, filename string, data []byte) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", dir, err)
	}

	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %q: %w", path, err)
	}

	return nil
}

// WriteFiles writes multiple files into dir at once.
// files is a map of filename -> content.
func WriteFiles(dir string, files map[string][]byte) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", dir, err)
	}

	for filename, data := range files {
		path := filepath.Join(dir, filename)
		if err := os.WriteFile(path, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %q: %w", path, err)
		}
	}

	return nil
}
