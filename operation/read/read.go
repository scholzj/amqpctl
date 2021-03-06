package read

import (
	"qpid.apache.org/amqp"
	"github.com/scholzj/amqpctl/mgmtlink"
	"github.com/scholzj/amqpctl/formatter"
	"bytes"
	"fmt"
	"errors"
)

func Read(link mgmtlink.MgmtLink, entityType string, attributeName string, attributeValue string) (output bytes.Buffer, err error) {
	reqProperties := make(map[string]interface{})
	reqProperties["operation"] = "READ"

	if attributeName == "identity" {
		reqProperties["identity"] = attributeValue
	} else if attributeName == "name" {
		reqProperties["name"] = attributeValue
		// Ready for WD11
		reqProperties["index"] = attributeName
		reqProperties["key"] = attributeValue
	} else {
		// Ready for WD11
		reqProperties["index"] = attributeName
		reqProperties["key"] = attributeValue
	}

	if entityType != "" {
		reqProperties["type"] = entityType
	}

	respProperties, respBody, err := link.Operation(reqProperties, nil)

	var statusCode int
	switch respProperties["statusCode"].(type) {
	case int32:
		statusCode = int(respProperties["statusCode"].(int32))
	case int64:
		statusCode = int(respProperties["statusCode"].(int64))
	}

	if err == nil {
		if statusCode == 200 {
			headers := []string{"ATTRIBUTE", "VALUE"}
			rows := parseReadResults(respBody)
			output = formatter.FormatPlainText(headers, rows)
		} else if statusCode == 400 {
			err = errors.New(fmt.Sprintf("Bad request: %v (%v)\n", respProperties["statusCode"], respProperties["statusDescription"]))
		} else if statusCode == 404 {
			err = errors.New(fmt.Sprintf("No manageable entities matching the request criteria found: %v (%v)\n", respProperties["statusCode"], respProperties["statusDescription"]))
		} else {
			err = errors.New(fmt.Sprintf("AMQP Management operation wasn't successfull: %v (%v)\n", respProperties["statusCode"], respProperties["statusDescription"]))
		}
	} else {
		err = errors.New(fmt.Sprintf("AMQP Management operation failed: %v\n", err.Error()))
	}

	return
}

func parseReadResults(body interface{}) (rows [][]string) {
	for attribute, value := range map[interface{}]interface{}(body.(amqp.Map)) {
		row := make([]string, 2)
		row[0] = attribute.(string)
		row[1] = fmt.Sprintf("%v", value)
		rows = append(rows, row)
	}

	return
}