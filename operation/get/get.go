package get

import (
	"fmt"
	"github.com/scholzj/amqpctl/mgmtlink"
	"errors"
)

func Get(link mgmtlink.MgmtLink, getOperation string, entityType string) (body interface{}, err error) {
	var reqProperties map[string]interface{}

	if entityType != "" {
		reqProperties = map[string]interface{}{"operation": getOperation, "entityType": entityType}
	} else {
		reqProperties = map[string]interface{}{"operation": getOperation}
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
			body = respBody
		} else {
			err = errors.New(fmt.Sprintf("AMQP Management operation wasn't successfull: %v (%v)\n", respProperties["statusCode"], respProperties["statusDescription"]))
			body = nil
		}
	} else {
		err = errors.New(fmt.Sprintf("AMQP Management operation failed: %v\n", err.Error()))
		body = nil
	}

	return
}