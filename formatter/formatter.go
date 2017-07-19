package formatter

import "bytes"

type Formatter interface {
	Format(headers []string, content [][]string) (bytes.Buffer)
}

type OutputFormat int

const (
	PLAINTEXT OutputFormat = iota
	JSON
	YAML
)
