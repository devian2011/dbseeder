package exporter

import (
	"dbseeder/internal/schema"
	"gopkg.in/yaml.v3"
	"io"
)

// ExportNotation export notation to writer
func ExportNotation(w io.Writer, s *schema.Schema) error {
	return yaml.NewEncoder(w).Encode(s.Databases)
}
