package configs

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// LanguageExtensionsData holds the raw byte content of the 'language_extensions.json' file.
// It is embedded into the binary at compile time using the go:embed directive.
//
//go:embed language_extensions.json
var LanguageExtensionsData []byte

var (
	// languageExtensions caches the parsed mapping of language names to their file extensions.
	languageExtensions map[string][]string
	// loadOnce ensures that the JSON unmarshaling logic is executed exactly once.
	loadOnce sync.Once
	// loadErr captures any error that occurs during the initial loading of the data.
	loadErr error
)

// LoadLanguageExtensions parses the embedded JSON data and returns a map where
// keys are lowercase language names (e.g., "go") and values are slices of
// associated file extensions (e.g., [".go"]).
//
// This function is thread-safe and uses a singleton pattern (sync.Once) to
// avoid redundant parsing on subsequent calls.
func LoadLanguageExtensions() (map[string][]string, error) {
	loadOnce.Do(func() {
		// langDefs defines the internal structure of the JSON entries for unmarshaling.
		var langDefs []struct {
			Name       string   `json:"name"`
			Type       string   `json:"type"`
			Extensions []string `json:"extensions"`
		}

		// Parse the embedded byte slice into the temporary slice of structs.
		if err := json.Unmarshal(LanguageExtensionsData, &langDefs); err != nil {
			loadErr = fmt.Errorf("failed to parse language extensions JSON: %w", err)
			return
		}

		// Initialize the global map and populate it with lowercase keys for case-insensitive lookup.
		languageExtensions = make(map[string][]string)
		for _, lang := range langDefs {
			langKey := strings.ToLower(lang.Name)
			languageExtensions[langKey] = lang.Extensions
		}
	})

	return languageExtensions, loadErr
}
