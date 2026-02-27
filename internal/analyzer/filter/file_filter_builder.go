package filter

import (
	"strings"

	"github.com/mhgffqwoer/git-fame/configs"
)

type FileFilterBuilder struct {
	includePatterns []pattern
	excludePatterns []pattern
	err             error
}

func NewFileFilterBuilder() *FileFilterBuilder {
	return &FileFilterBuilder{}
}

func (f *FileFilterBuilder) WithExtensions(patterns string) *FileFilterBuilder {
	if f.err != nil {
		return f
	}
	patterns = strings.ReplaceAll(patterns, "'", "")
	list := strings.SplitSeq(patterns, ",")
	for pattern := range list {
		if pattern == "" {
			continue
		}
		f.includePatterns = append(f.includePatterns, newFilterPattern(pattern))
	}
	return f
}

func (f *FileFilterBuilder) WithLanguages(languages string) *FileFilterBuilder {
	if f.err != nil {
		return f
	}
	langToExt, err := configs.LoadLanguageExtensions()
	if err != nil {
		f.err = err
		return f
	}
	languages = strings.ReplaceAll(languages, "'", "")
	list := strings.SplitSeq(languages, ",")
	for lang := range list {
		langLower := strings.ToLower(lang)
		if extensions, exists := langToExt[langLower]; exists {
			for _, ext := range extensions {
				f.includePatterns = append(f.includePatterns, newFilterPattern(ext))
			}
		}
	}
	return f
}

func (f *FileFilterBuilder) WithExcludes(patterns string) *FileFilterBuilder {
	if f.err != nil {
		return f
	}
	patterns = strings.ReplaceAll(patterns, "'", "")
	list := strings.SplitSeq(patterns, ",")
	for pattern := range list {
		if pattern == "" {
			continue
		}
		f.excludePatterns = append(f.excludePatterns, newFilterPattern(pattern))
	}
	return f
}

func (f *FileFilterBuilder) WithRestrictTo(patterns string) *FileFilterBuilder {
	if f.err != nil {
		return f
	}
	patterns = strings.ReplaceAll(patterns, "'", "")
	list := strings.SplitSeq(patterns, ",")
	for pattern := range list {
		if pattern == "" {
			continue
		}
		f.includePatterns = append(f.includePatterns, newFilterPattern(pattern))
	}
	return f
}

func (f *FileFilterBuilder) Build() (*FileFilter, error) {
	if f.err != nil {
		return nil, f.err
	}
	return newFileFilter(f.includePatterns, f.excludePatterns), nil
}
