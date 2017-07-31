package query

import (
	"qpid.apache.org/amqp"
	"strings"
	"github.com/scholzj/amqpctl/mgmtlink"
	"github.com/scholzj/amqpctl/formatter"
	"bytes"
	"fmt"
	"errors"
)

func Query(link mgmtlink.MgmtLink, entityType string, attributes []string) (output bytes.Buffer, err error) {
	var reqProperties map[string]interface{}

	if entityType != "" {
		reqProperties = map[string]interface{}{"operation": "QUERY", "entityType": entityType}
	} else {
		reqProperties = map[string]interface{}{"operation": "QUERY"}
	}

	reqBody := map[interface{}]interface{}{"attributeNames": attributes}

	respProperties, respBody, err := link.Operation(reqProperties, reqBody)

	var statusCode int
	switch respProperties["statusCode"].(type) {
	case int32:
		statusCode = int(respProperties["statusCode"].(int32))
	case int64:
		statusCode = int(respProperties["statusCode"].(int64))
	}

	if err == nil {
		if statusCode == 200 {
			headers := parseQueryHeaders(respBody)
			rows := parseQueryResults(respBody)
			output = formatter.FormatPlainText(headers, rows)
		} else {
			err = errors.New(fmt.Sprintf("AMQP Management operation wasn't successfull: %v (%v)\n", respProperties["statusCode"], respProperties["statusDescription"]))
		}
	} else {
		err = errors.New(fmt.Sprintf("AMQP Management operation failed: %v\n", err.Error()))
	}

	return
}

func parseQueryHeaders(body interface{}) (headers []string) {
	attributeNames := map[interface{}]interface{}(body.(amqp.Map))["attributeNames"]

	for _, attribute := range []interface{}(attributeNames.(amqp.List)) {
		headers = append(headers, strings.ToUpper(attribute.(string)))
	}

	return
}

func parseQueryResults(body interface{}) (rows [][]string) {
	results := map[interface{}]interface{}(body.(amqp.Map))["results"]

	for _, attributeList := range []interface{}(results.(amqp.List)) {
		row := []string{}
		for _, attribute := range []interface{}(attributeList.(amqp.List)) {
			row = append(row, fmt.Sprintf("%v", attribute))
		}

		rows = append(rows, row)
	}

	return
}