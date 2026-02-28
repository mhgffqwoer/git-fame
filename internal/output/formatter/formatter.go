package formatter

import (
	"fmt"

	metrics "github.com/mhgffqwoer/git-fame/internal/metrics"
	format_types "github.com/mhgffqwoer/git-fame/internal/output/formatter/format_types"
)

// FormatType defines the interface for rendering formatted output.
type FormatType interface {
	Render(data []*metrics.Data) (string, error)
}

// Formatter represents a formatter for rendering formatted output.
type Formatter struct {
	formatter FormatType // The specific format type used for rendering (e.g., tabular, CSV, JSON).
}

// NewFormatter creates a new Formatter instance based on the provided formatType.
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

// Format renders the formatted output based on the provided data.
func (f *Formatter) Format(data []*metrics.Data) (string, error) {
	return f.formatter.Render(data)
}
