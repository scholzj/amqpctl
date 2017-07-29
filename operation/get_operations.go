package operation

import (
	"qpid.apache.org/amqp"
	"strings"
	"github.com/scholzj/amqpctl/mgmtlink"
	"github.com/scholzj/amqpctl/formatter"
	"bytes"
)

func GetOperations(link mgmtlink.MgmtLink, entityType string) (output bytes.Buffer, err error) {
	body, err := Get(link, "GET-OPERATIONS", entityType)

	if err == nil {
		rows := parseGetOperationsResults(body)
		output = formatter.FormatPlainText([]string{"TYPE", "OPERATION", "ARGUMENT"}, rows)
	}

	return
}

func parseGetOperationsResults(body interface{}) (rows [][]string) {
	for entityType, operations := range map[interface{}]interface{}(body.(amqp.Map)) {
		firstLine := true

		switch operations.(type) {
		case amqp.Map:
			for operation, argumentsList := range map[interface{}]interface{}(operations.(amqp.Map)) {
				argumentsArray := make([]string, len([]interface{}(argumentsList.(amqp.List))))

				for i, argument := range []interface{}(argumentsList.(amqp.List)) {
					argumentsArray[i] = argument.(string)
				}

				var entity interface{} = ""
				if firstLine {
					entity = entityType
					firstLine = false
				}

				row := make([]string, 3)
				row[0] = entity.(string)
				row[1] = operation.(string)
				row[2] = strings.Join(argumentsArray, ", ")
				rows = append(rows, row)
			}
		case amqp.List:
			for _, operation := range []interface{}(operations.(amqp.List)) {
				var entity interface{} = ""
				if firstLine {
					entity = entityType
					firstLine = false
				}

				row := make([]string, 3)
				row[0] = entity.(string)
				row[1] = operation.(string)
				row[2] = ""
				rows = append(rows, row)
			}
		}
	}

	return
}