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

func Query(args []string) {
	usage := `Usage:
  amqpctl query [<entityType>] [<attributes>...]

Options:
  -h --help   Show this screen.

Description:
  Query selected attributes of Management entities
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
		reqProperties = map[string]interface{}{"operation": "QUERY", "entityType": arguments["<entityType>"]}
	} else {
		reqProperties = map[string]interface{}{"operation": "QUERY"}
	}

	var reqBody map[interface{}]interface{}

	if len(arguments["<attributes>"].([]string)) > 0 {
		if len(arguments["<attributes>"].([]string)) == 1 && arguments["<attributes>"].([]string)[0] == "%" {
			reqBody = map[interface{}]interface{}{"attributeNames": []string{}}
		} else {
			reqBody = map[interface{}]interface{}{"attributeNames": arguments["<attributes>"].([]string)}
		}

	} else {
		reqBody = map[interface{}]interface{}{"attributeNames": []interface{}{"name", "type", "identity"}}
	}

	respProperties, respBody, err := link.Operation(reqProperties, reqBody)

	if err == nil {
		if respProperties["statusCode"].(int32) == 200 {
			printQueryResults(respProperties, respBody)
		} else {
			fmt.Printf("ERROR %v: %v\n", respProperties["statusCode"], respProperties["statusDescription"])
			os.Exit(1)
		}
	} else {
		fmt.Printf("Ups, something went wrong: %v\n", err.Error())
		os.Exit(1)
	}
}

func printQueryResults(properties map[string]interface{}, body interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 3, ' ', 0)

	// Print dynamic header
	attributeNames := map[interface{}]interface{}(body.(amqp.Map))["attributeNames"]

	for _, attribute := range []interface{}(attributeNames.(amqp.List)) {
		fmt.Fprintf(w, "%v\t", attribute.(string))
	}

	fmt.Fprint(w, "\n")

	// Print content
	results := map[interface{}]interface{}(body.(amqp.Map))["results"]

	for _, attributeList := range []interface{}(results.(amqp.List)) {
		for _, attribute := range []interface{}(attributeList.(amqp.List)) {
			fmt.Fprintf(w, "%v\t", attribute)
		}

		fmt.Fprint(w, "\n")
	}

	w.Flush()
}
