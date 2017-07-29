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
		rows := parseResults(body)
		output = formatter.FormatPlainText([]string{"TYPE", "PARENT"}, rows)
	}

	return
}

/*func GetTypes(link mgmtlink.MgmtLink, entityType string) (output bytes.Buffer, err error) {
	var reqProperties map[string]interface{}

	if entityType != "" {
		reqProperties = map[string]interface{}{"operation": "GET-TYPES", "entityType": entityType}
	} else {
		reqProperties = map[string]interface{}{"operation": "GET-TYPES"}
	}

	respProperties, respBody, err := link.Operation(reqProperties, nil)

	if err == nil {
		if respProperties["statusCode"].(int64) == 200 {
			rows := parseResults(respBody)
			output = formatter.FormatPlainText([]string{"TYPE", "PARENT"}, rows)
		} else {
			err = errors.New(fmt.Sprintf("AMQP Management operation wasn't successfull: %v (%v)\n", respProperties["statusCode"], respProperties["statusDescription"]))
		}
	} else {
		err = errors.New(fmt.Sprintf("AMQP Management operation failed: %v\n", err.Error()))
	}

	return
}*/

func parseResults(body interface{}) (rows [][]string) {
	rows = make([][]string, len(map[interface{}]interface{}(body.(amqp.Map))))

	i := 0
	for entitytype, extends := range map[interface{}]interface{}(body.(amqp.Map)) {
		parents := make([]string, len([]interface{}(extends.(amqp.List))))
		for j, parent := range []interface{}(extends.(amqp.List)) {
			parents[j] = parent.(string)
		}

		row := make([]string, 2)
		row[0] = entitytype.(string)
		row[1] = strings.Join(parents, ", ")
		rows[i] = row
		i++
	}

	return
}