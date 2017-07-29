package operation

import (
	"qpid.apache.org/amqp"
	"strings"
	"github.com/scholzj/amqpctl/mgmtlink"
	"github.com/scholzj/amqpctl/formatter"
	"bytes"
)

func GetTypes(link mgmtlink.MgmtLink, entityType string) (output bytes.Buffer, err error) {
	body, err := Get(link, "GET-TYPES", entityType)

	if err == nil {
		rows := parseGetTypesResults(body)
		output = formatter.FormatPlainText([]string{"TYPE", "PARENT"}, rows)
	}

	return
}

func parseGetTypesResults(body interface{}) (rows [][]string) {
	for entitytype, extends := range map[interface{}]interface{}(body.(amqp.Map)) {
		parents := make([]string, len([]interface{}(extends.(amqp.List))))
		for j, parent := range []interface{}(extends.(amqp.List)) {
			parents[j] = parent.(string)
		}

		row := make([]string, 2)
		row[0] = entitytype.(string)
		row[1] = strings.Join(parents, ", ")
		rows = append(rows, row)
	}

	return
}