package filter

import (
	"strings"

	"github.com/mhgffqwoer/git-fame/configs"
)

// FileFilterBuilder collects various filtering criteria and builds a FileFilter.
type FileFilterBuilder struct {
	includePatterns []pattern // Patterns that files must match to be included.
	excludePatterns []pattern // Patterns that files must NOT match to be included.
	err             error     // Sticky error to simplify builder method chaining.
}

// NewFileFilterBuilder initializes an empty builder.
func NewFileFilterBuilder() *FileFilterBuilder {
	return &FileFilterBuilder{}
}

// WithExtensions adds a comma-separated list of file extensions (e.g., ".go,.md") to the inclusion list.
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

// WithLanguages resolves language names (e.g., "Go") to extensions using the embedded config
// and adds them to the inclusion list.
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

// WithExcludes adds glob patterns to the exclusion list. Any file matching these will be ignored.
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

// WithRestrictTo adds patterns that files must match to be included.
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

// Build validates the builder state and returns the final FileFilter.
func (f *FileFilterBuilder) Build() (*FileFilter, error) {
	if f.err != nil {
		return nil, f.err
	}
	return newFileFilter(f.includePatterns, f.excludePatterns), nil
}
