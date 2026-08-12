package format_types

import (
	"bytes"
	"encoding/csv"
	"strconv"

	metrics "github.com/mhgffqwoer/git-fame/internal/metrics"
)

// CSVFormatType renders data in CSV format.
type CSVFormatType struct{}

// NewCSVFormatType creates a new CSVFormatType instance.
func NewCSVFormatType() *CSVFormatType {
	return &CSVFormatType{}
}

// Render renders data in CSV format.
func (t *CSVFormatType) Render(data []*metrics.Data) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	err := writer.Write([]string{"Name", "Lines", "Commits", "Files"})
	if err != nil {
		return "", err
	}
	for _, d := range data {
		record := []string{
			d.Name(),
			strconv.Itoa(d.LinesCount()),
			strconv.Itoa(d.CommitsCount()),
			strconv.Itoa(d.FilesCount()),
		}
		if err := writer.Write(record); err != nil {
			return "", err
		}
	}
	writer.Flush()
	return buf.String(), nil
}
