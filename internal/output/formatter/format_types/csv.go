package format_types

import (
	"bytes"
	"encoding/csv"
	"strconv"

	metrics "gitlab.com/slon/shad-go/gitfame/internal/metrics"
)

type CSVFormatType struct{}

func NewCSVFormatType() *CSVFormatType {
	return &CSVFormatType{}
}

func (t *CSVFormatType) Render(data []*metrics.Data) (string, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write([]string{"Name", "Lines", "Commits", "Files"}); err != nil {
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
