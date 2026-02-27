package entity

import "path/filepath"

type File struct {
	Path string
}

func NewFile(path string) *File {
	return &File{
		Path: path,
	}
}

func (f File) Directory() string {
	return filepath.Dir(f.Path)
}

func (f File) Name() string {
	return filepath.Base(f.Path)
}

func (f File) Extension() string {
	return filepath.Ext(f.Path)
}

func (f File) String() string {
	return f.Path
}
