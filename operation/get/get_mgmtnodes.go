package get

import (
	"qpid.apache.org/amqp"
	"github.com/scholzj/amqpctl/mgmtlink"
	"github.com/scholzj/amqpctl/formatter"
	"bytes"
)

func GetMgmtNodes(link mgmtlink.MgmtLink, entityType string) (output bytes.Buffer, err error) {
	body, err := Get(link, "GET-MGMT-NODES", entityType)

	if err == nil {
		rows := parseGetMgmtNodesResults(body)
		output = formatter.FormatPlainText([]string{"NODE"}, rows)
	}

	return
}

func parseGetMgmtNodesResults(body interface{}) (rows [][]string) {
	for _, address := range []interface{}(body.(amqp.List)) {
		row := make([]string, 1)
		row[0] = address.(string)
		rows = append(rows, row)
	}

	return
}