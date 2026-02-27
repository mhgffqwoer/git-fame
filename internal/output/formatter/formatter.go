package formatter

import (
	"fmt"

	metrics "github.com/mhgffqwoer/git-fame/internal/metrics"
	format_types "github.com/mhgffqwoer/git-fame/internal/output/formatter/format_types"
)

type FormatType interface {
	Render(data []*metrics.Data) (string, error)
}

type Formatter struct {
	formatter FormatType
}

func NewFormatter(formatType string) (*Formatter, error) {
	switch formatType {
	case "tabular":
		return &Formatter{
			formatter: format_types.NewTabularFormatType(),
		}, nil
	case "csv":
		return &Formatter{
			formatter: format_types.NewCSVFormatType(),
		}, nil
	case "json":
		return &Formatter{
			formatter: format_types.NewJSONFormatType(),
		}, nil
	case "json-lines":
		return &Formatter{
			formatter: format_types.NewJSONLinesFormatType(),
		}, nil
	default:
		return nil, fmt.Errorf("error")
	}
}

func (f *Formatter) Format(data []*metrics.Data) (string, error) {
	return f.formatter.Render(data)
}
