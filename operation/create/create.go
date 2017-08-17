package create

import (
	"qpid.apache.org/amqp"
	"github.com/scholzj/amqpctl/mgmtlink"
	"github.com/scholzj/amqpctl/formatter"
	"bytes"
	"fmt"
	"errors"
)

func Create(link mgmtlink.MgmtLink, createMap map[interface{}]interface{}) (output bytes.Buffer, err error) {
	var reqProperties map[string]interface{}
	var reqBody map[interface{}]interface{}

	reqProperties = map[string]interface{}{"operation": "CREATE"}

	reqBody = createMap

	respProperties, respBody, err := link.Operation(reqProperties, reqBody)

	var statusCode int
	switch respProperties["statusCode"].(type) {
	case int32:
		statusCode = int(respProperties["statusCode"].(int32))
	case int64:
		statusCode = int(respProperties["statusCode"].(int64))
	}

	if err == nil {
		if statusCode == 201 {
			headers := []string{"ATTRIBUTE", "VALUE"}
			rows := parseCreateResults(respBody)
			output = formatter.FormatPlainText(headers, rows)
		} else if statusCode == 400 {
			err = errors.New(fmt.Sprintf("Specified index is not supported: %v (%v)\n", respProperties["statusCode"], respProperties["statusDescription"]))
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

func parseCreateResults(body interface{}) (rows [][]string) {
	for attribute, value := range map[interface{}]interface{}(body.(amqp.Map)) {
		row := make([]string, 2)
		row[0] = attribute.(string)
		row[1] = fmt.Sprintf("%v", value)
		rows = append(rows, row)
	}

	return
}