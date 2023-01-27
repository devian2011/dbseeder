package exporter

import (
	"dbseeder/internal/schema"
	"gopkg.in/yaml.v3"
	"io"
)

func ExportNotation(w io.Writer, s *schema.Schema) error {
	return yaml.NewEncoder(w).Encode(s.Databases)
}
