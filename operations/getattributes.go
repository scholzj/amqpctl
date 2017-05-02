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

func GetAttributes(args []string) {
	usage := `Usage:
  amqpctl getattributes [<entityType>]

Options:
  -h --help   Show this screen.

Description:
  Get list of attributes implemented supported by given manageable entity type.
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
		reqProperties = map[string]interface{}{"operation": "GET-ATTRIBUTES", "entityType": arguments["<entityType>"]}
	} else {
		reqProperties = map[string]interface{}{"operation": "GET-ATTRIBUTES"}
	}

	respProperties, respBody, err := link.Operation(reqProperties, nil)

	if err == nil {
		printAttributes(respProperties, respBody)
	} else {
		fmt.Printf("Ups, something went wrong: %v\n", err.Error())
		os.Exit(1)
	}
}

func printAttributes(properties map[string]interface{}, body interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 3, ' ', 0)
	fmt.Fprint(w, "TYPE\tATTRIBUTES\t\n")

	for entitytype, attributesList := range map[interface{}]interface{}(body.(amqp.Map)) {
		attributesArray := make([]string, len([]interface{}(attributesList.(amqp.List))))
		for i, attribute := range []interface{}(attributesList.(amqp.List)) {
			attributesArray[i] = attribute.(string)
		}


		fmt.Fprintf(w, "%v\t%v\t\n", entitytype, strings.Join(attributesArray, ", "))
	}

	w.Flush()
}
