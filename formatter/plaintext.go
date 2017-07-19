package formatter

import (
	"text/tabwriter"
	"fmt"
	"strings"
	"bytes"
)

type PlainTextFormatter struct {

}

func (f *PlainTextFormatter) Format(headers []string, rows [][]string) (bytes.Buffer) {
	buf := new(bytes.Buffer)
	w := tabwriter.NewWriter(buf, 10, 4, 3, ' ', 0)

	// Print header
	fmt.Fprintf(w, "%s\n", strings.Join(headers, "\t"))

	// PrintbBody
	for _, row := range rows {
		fmt.Fprintf(w, "%s\n", strings.Join(row, "\t"))
	}

	w.Flush()

	return *buf
}