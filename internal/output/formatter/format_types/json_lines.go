package format_types

import (
	"encoding/json"
	"strings"

	metrics "github.com/mhgffqwoer/git-fame/internal/metrics"
)

type JSONLinesFormatType struct{}

func NewJSONLinesFormatType() *JSONLinesFormatType {
	return &JSONLinesFormatType{}
}

func (t *JSONLinesFormatType) Render(data []*metrics.Data) (string, error) {
	type outData struct {
		Name    string `json:"name"`
		Lines   int    `json:"lines"`
		Commits int    `json:"commits"`
		Files   int    `json:"files"`
	}
	var b strings.Builder
	for i, d := range data {
		record := outData{
			Name:    d.Name(),
			Lines:   d.LinesCount(),
			Commits: d.CommitsCount(),
			Files:   d.FilesCount(),
		}
		jsonData, err := json.Marshal(record)
		if err != nil {
			return "", err
		}
		b.Write(jsonData)
		if i < len(data)-1 {
			b.WriteByte('\n')
		}
	}
	return b.String(), nil
}
