package format_types

import (
	"encoding/json"

	metrics "gitlab.com/slon/shad-go/gitfame/internal/metrics"
)

type JSONFormatType struct{}

func NewJSONFormatType() *JSONFormatType {
	return &JSONFormatType{}
}

func (t *JSONFormatType) Render(data []*metrics.Data) (string, error) {
	type outData struct {
		Name    string `json:"name"`
		Lines   int    `json:"lines"`
		Commits int    `json:"commits"`
		Files   int    `json:"files"`
	}
	result := make([]outData, 0, len(data))
	for _, d := range data {
		result = append(result, outData{
			Name:    d.Name(),
			Lines:   d.LinesCount(),
			Commits: d.CommitsCount(),
			Files:   d.FilesCount(),
		})
	}
	jsonBytes, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}
