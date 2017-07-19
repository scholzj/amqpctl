package operation

import (
	"fmt"
	"os"
	"text/tabwriter"
	"qpid.apache.org/amqp"
	"strings"
	"github.com/scholzj/amqpctl/mgmtlink"
	"github.com/scholzj/amqpctl/formatter"
)

func getTypes(link mgmtlink.MgmtLink, args []string, output formatter.OutputFormat) {
	var reqProperties map[string]interface{}

	if len(args) > 0 {
		reqProperties = map[string]interface{}{"operation": "GET-TYPES", "entityType": args[0]}
	} else {
		reqProperties = map[string]interface{}{"operation": "GET-TYPES"}
	}

	respProperties, respBody, err := link.Operation(reqProperties, nil)

	if err == nil {
		if respProperties["statusCode"].(int64) == 200 {
			printTypes(respProperties, respBody)
		} else {
			fmt.Printf("AMQP Management operation wsn't successfull: %v (%v)\n", respProperties["statusCode"], respProperties["statusDescription"])
			os.Exit(1)
		}
	} else {
		fmt.Printf("AMQP Management operation failed: %v\n", err.Error())
		os.Exit(1)
	}

	//jsonString, _ := json.Marshal(respBody)

	/*if output == formatter.YAML {

	} else if output == formatter.JSON {
		return nil
	} else {
		f := formatter.PlainTextFormatter{}
		f.Format()
	}*/
}

func printTypes(properties map[string]interface{}, body interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 3, ' ', 0)
	fmt.Fprint(w, "TYPE\tPARENTS\t\n")

	for entitytype, extends := range map[interface{}]interface{}(body.(amqp.Map)) {
		parents := make([]string, len([]interface{}(extends.(amqp.List))))
		for i, parent := range []interface{}(extends.(amqp.List)) {
			parents[i] = parent.(string)
		}


		fmt.Fprintf(w, "%v\t%v\t\n", entitytype, strings.Join(parents, ", "))
	}

	w.Flush()
}