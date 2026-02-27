package format_types

import (
	"fmt"
	"strings"
	"text/tabwriter"

	metrics "gitlab.com/slon/shad-go/gitfame/internal/metrics"
)

type TabularFormatType struct{}

func NewTabularFormatType() *TabularFormatType {
	return &TabularFormatType{}
}

func (t *TabularFormatType) Render(data []*metrics.Data) (string, error) {
	var b strings.Builder
	writer := tabwriter.NewWriter(&b, 0, 0, 1, ' ', 0)
	if _, err := writer.Write([]byte("Name\tLines\tCommits\tFiles\n")); err != nil {
		return "", err
	}
	for _, d := range data {
		record := []byte(d.Name() + "\t" +
			fmt.Sprint(d.LinesCount()) + "\t" +
			fmt.Sprint(d.CommitsCount()) + "\t" +
			fmt.Sprint(d.FilesCount()) + "\n")
		if _, err := writer.Write(record); err != nil {
			return "", err
		}
	}
	if err := writer.Flush(); err != nil {
		return "", err
	}
	return b.String(), nil
}
