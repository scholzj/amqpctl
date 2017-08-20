package delete

import (
	"github.com/scholzj/amqpctl/mgmtlink"
	"fmt"
	"errors"
)

func Delete(link mgmtlink.MgmtLink, entityType string, attributeName string, attributeValue string) (err error) {
	reqProperties := make(map[string]interface{})
	reqProperties["operation"] = "DELETE"

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

	respProperties, _, err := link.Operation(reqProperties, nil)

	var statusCode int
	switch respProperties["statusCode"].(type) {
	case int32:
		statusCode = int(respProperties["statusCode"].(int32))
	case int64:
		statusCode = int(respProperties["statusCode"].(int64))
	}

	if err == nil {
		if statusCode == 204 {
			err = nil
		} else if statusCode == 400 {
			err = errors.New(fmt.Sprintf("Bad Request: %v (%v)\n", respProperties["statusCode"], respProperties["statusDescription"]))
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
