package configs

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

//go:embed language_extensions.json
var LanguageExtensionsData []byte

var (
	languageExtensions map[string][]string
	loadOnce           sync.Once
	loadErr            error
)

func LoadLanguageExtensions() (map[string][]string, error) {
	loadOnce.Do(func() {
		var langDefs []struct {
			Name       string   `json:"name"`
			Type       string   `json:"type"`
			Extensions []string `json:"extensions"`
		}

		if err := json.Unmarshal(LanguageExtensionsData, &langDefs); err != nil {
			loadErr = fmt.Errorf("failed to parse language extensions JSON: %w", err)
			return
		}

		languageExtensions = make(map[string][]string)
		for _, lang := range langDefs {
			langKey := strings.ToLower(lang.Name)
			languageExtensions[langKey] = lang.Extensions
		}
	})

	return languageExtensions, loadErr
}
