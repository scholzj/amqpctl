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

func GetAnnotations(args []string) {
	usage := `Usage:
  amqpctl getannotations [<entityType>]

Options:
  -h --help   Show this screen.

Description:
  Get list of annotations implemented by the manageable entity types.
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
		reqProperties = map[string]interface{}{"operation": "GET-ANNOTATIONS", "entityType": arguments["<entityType>"]}
	} else {
		reqProperties = map[string]interface{}{"operation": "GET-ANNOTATIONS"}
	}

	respProperties, respBody, err := link.Operation(reqProperties, nil)

	if err == nil {
		printAnnotations(respProperties, respBody)
	} else {
		fmt.Printf("Ups, something went wrong: %v\n", err.Error())
		os.Exit(1)
	}
}

func printAnnotations(properties map[string]interface{}, body interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 3, ' ', 0)
	fmt.Fprint(w, "TYPE\tANNOTATIONS\t\n")

	for entitytype, annotationsList := range map[interface{}]interface{}(body.(amqp.Map)) {
		annotationsArray := make([]string, len([]interface{}(annotationsList.(amqp.List))))
		for i, annotation := range []interface{}(annotationsList.(amqp.List)) {
			annotationsArray[i] = annotation.(string)
		}


		fmt.Fprintf(w, "%v\t%v\t\n", entitytype, strings.Join(annotationsArray, ", "))
	}

	w.Flush()
}
