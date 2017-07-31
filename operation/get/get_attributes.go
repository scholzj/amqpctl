package get

import (
	"qpid.apache.org/amqp"
	"strings"
	"github.com/scholzj/amqpctl/mgmtlink"
	"github.com/scholzj/amqpctl/formatter"
	"bytes"
)

func GetAttributes(link mgmtlink.MgmtLink, entityType string) (output bytes.Buffer, err error) {
	body, err := Get(link, "GET-ATTRIBUTES", entityType)

	if err == nil {
		rows := parseGetAttributesResults(body)
		output = formatter.FormatPlainText([]string{"TYPE", "ATTRIBUTES"}, rows)
	}

	return
}

func parseGetAttributesResults(body interface{}) (rows [][]string) {
	for entitytype, attributesList := range map[interface{}]interface{}(body.(amqp.Map)) {
		attributesArray := make([]string, len([]interface{}(attributesList.(amqp.List))))
		for i, attribute := range []interface{}(attributesList.(amqp.List)) {
			attributesArray[i] = attribute.(string)
		}

		row := make([]string, 2)
		row[0] = entitytype.(string)
		row[1] = strings.Join(attributesArray, ", ")
		rows = append(rows, row)
	}

	return
}