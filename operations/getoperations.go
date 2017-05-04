package operations

import (
	"github.com/docopt/docopt-go"
	"fmt"
	"strings"
	"os"
	"text/tabwriter"
	"../utils"
	"qpid.apache.org/amqp"
)

func GetOperations(args []string, link utils.MgmtLink) {
	usage := `Usage:
  amqpctl getoperations [<entityType>]

Options:
  -h --help   Show this screen.

Description:
  Get list of operations implemented supported by given manageable entity type.
`
	arguments, err := docopt.Parse(usage, args, true, "", false, false)
	if err != nil {
		fmt.Printf("Invalid option: 'amqpctl %s'. Use flag '--help' to read about a specific subcommand.\n", strings.Join(args, " "))
		os.Exit(1)
	}
	if len(arguments) == 0 {
		return
	}

	err = link.Connect()
	if err != nil {
		fmt.Printf("Ups, something went wrong ... %v\n", err.Error())
		os.Exit(1)
	}

	defer link.Close()

	var reqProperties map[string]interface{}

	if arguments["<entityType>"] != nil {
		reqProperties = map[string]interface{}{"operation": "GET-OPERATIONS", "entityType": arguments["<entityType>"]}
	} else {
		reqProperties = map[string]interface{}{"operation": "GET-OPERATIONS"}
	}

	respProperties, respBody, err := link.Operation(reqProperties, nil)

	if err == nil {
		if respProperties["statusCode"].(int32) == 200 {
			printOperations(respProperties, respBody)
		} else {
			fmt.Printf("Ups, something went wrong ... %v: %v\n", respProperties["statusCode"], respProperties["statusDescription"])
			os.Exit(1)
		}
	} else {
		fmt.Printf("Ups, something went wrong ... %v\n", err.Error())
		os.Exit(1)
	}
}

func printOperations(properties map[string]interface{}, body interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 3, ' ', 0)
	fmt.Fprint(w, "TYPE\tOPERATION\tARGUMENTS\t\n")

	for entityType, operations := range map[interface{}]interface{}(body.(amqp.Map)) {
		firstLine := true

		switch operations.(type) {
		case amqp.Map:
			for operation, argumentsList := range map[interface{}]interface{}(operations.(amqp.Map)) {
				argumentsArray := make([]string, len([]interface{}(argumentsList.(amqp.List))))

				for i, argument := range []interface{}(argumentsList.(amqp.List)) {
					argumentsArray[i] = argument.(string)
				}

				var entity interface{} = ""
				if firstLine {
					entity = entityType
					firstLine = false
				}

				fmt.Fprintf(w, "%v\t%v\t%v\t\n", entity, operation, strings.Join(argumentsArray, ", "))
			}
		case amqp.List:
			for _, operation := range []interface{}(operations.(amqp.List)) {
				var entity interface{} = ""
				if firstLine {
					entity = entityType
					firstLine = false
				}

				fmt.Fprintf(w, "%v\t%v\t%v\t\n", entity, operation, "")
			}
			/*for _, operationsList := range []interface{}(operations.(amqp.List)) {
				for _, operation := range []interface{}(operationsList) {
					var entity interface{} = ""
					if firstLine {
						entity = entityType
						firstLine = false
					}

					fmt.Fprintf(w, "%v\t%v\t%v\t\n", entity, operation, "")
				}
			}*/
		}
	}

	w.Flush()
}
