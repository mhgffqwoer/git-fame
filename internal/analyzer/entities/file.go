package entity

import "path/filepath"

// File represents a single file within the git repository being analyzed.
type File struct {
	Path string // The relative path to the file from the repository root.
}

// NewFile creates a new File instance with the given path.
func NewFile(path string) *File {
	return &File{
		Path: path,
	}
}

// Directory returns the directory portion of the file path.
func (f File) Directory() string {
	return filepath.Dir(f.Path)
}

// Name returns the base name (filename) of the file.
func (f File) Name() string {
	return filepath.Base(f.Path)
}

// Extension returns the file extension (e.g., ".go").
func (f File) Extension() string {
	return filepath.Ext(f.Path)
}

// String returns the full path of the file as a string.
func (f File) String() string {
	return f.Path
}
