package operations

import (
	"github.com/docopt/docopt-go"
	"fmt"
	"strings"
	"os"
	"text/tabwriter"
	"./utils"
	"qpid.apache.org/amqp"
)

func GetTypes(args []string) {
	usage := `Usage:
  amqpctl gettypes [<entityType>]

Options:
  -h --help   Show this screen.

Description:
  Get list of supported manageable entity types.
`
	arguments, err := docopt.Parse(usage, args, true, "", false, false)
	if err != nil {
		fmt.Printf("Invalid option: 'amqpctl %s'. Use flag '--help' to read about a specific subcommand.\n", strings.Join(args, " "))
		os.Exit(1)
	}
	if len(arguments) == 0 {
		return
	}

	link := utils.MgmtLink{}
	err = link.Connect()
	if err != nil {
		fmt.Printf("Ups, something went wrong: %v\n", err.Error())
		os.Exit(1)
	}

	defer link.Close()

	var reqProperties map[string]interface{}

	if arguments["<entityType>"] != nil {
		reqProperties = map[string]interface{}{"operation": "GET-TYPES", "entityType": arguments["<entityType>"]}
	} else {
		reqProperties = map[string]interface{}{"operation": "GET-TYPES"}
	}

	respProperties, respBody, err := link.Operation(reqProperties, nil)

	if err == nil {
		if respProperties["statusCode"].(int32) == 200 {
			printTypes(respProperties, respBody)
		} else {
			fmt.Printf("ERROR %v: %v\n", respProperties["statusCode"], respProperties["statusDescription"])
			os.Exit(1)
		}
	} else {
		fmt.Printf("Ups, something went wrong: %v\n", err.Error())
		os.Exit(1)
	}
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
