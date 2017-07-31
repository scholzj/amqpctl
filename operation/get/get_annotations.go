package get

import (
	"qpid.apache.org/amqp"
	"strings"
	"github.com/scholzj/amqpctl/mgmtlink"
	"github.com/scholzj/amqpctl/formatter"
	"bytes"
)

func GetAnnotations(link mgmtlink.MgmtLink, entityType string) (output bytes.Buffer, err error) {
	body, err := Get(link, "GET-ANNOTATIONS", entityType)

	if err == nil {
		rows := parseGetAnnotationsResults(body)
		output = formatter.FormatPlainText([]string{"TYPE", "ANNOTATIONS"}, rows)
	}

	return
}

func parseGetAnnotationsResults(body interface{}) (rows [][]string) {
	for entitytype, annotationsList := range map[interface{}]interface{}(body.(amqp.Map)) {
		annotationsArray := make([]string, len([]interface{}(annotationsList.(amqp.List))))
		for i, annotation := range []interface{}(annotationsList.(amqp.List)) {
			annotationsArray[i] = annotation.(string)
		}

		row := make([]string, 2)
		row[0] = entitytype.(string)
		row[1] = strings.Join(annotationsArray, ", ")
		rows = append(rows, row)
	}

	return
}