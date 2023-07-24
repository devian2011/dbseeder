package exporter

import (
	"io"

	"gopkg.in/yaml.v3"

	"dbseeder/internal/schema"
)

// ExportNotation export notation to writer
func ExportNotation(w io.Writer, s *schema.Schema) error {
	return yaml.NewEncoder(w).Encode(s.Databases)
}
