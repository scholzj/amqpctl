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

func GetMgmtNodes(args []string) {
	usage := `Usage:
  amqpctl getmgmtnodes

Options:
  -h --help   Show this screen.

Description:
  Get list of addresses of other management nodes which this management node is aware of.
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
	reqProperties = map[string]interface{}{"operation": "GET-MGMT-NODES"}
	respProperties, respBody, err := link.Operation(reqProperties, nil)

	if err == nil {
		if respProperties["statusCode"].(int32) == 200 {
			printMgmtNodes(respProperties, respBody)
		} else {
			fmt.Printf("ERROR %v: %v\n", respProperties["statusCode"], respProperties["statusDescription"])
			os.Exit(1)
		}
	} else {
		fmt.Printf("Ups, something went wrong: %v\n", err.Error())
		os.Exit(1)
	}
}

func printMgmtNodes(properties map[string]interface{}, body interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 10, 4, 3, ' ', 0)
	fmt.Fprint(w, "MGMTNODE\t\n")

	for _, address := range []interface{}(body.(amqp.List)) {
		fmt.Fprintf(w, "%v\t\n", address)
	}

	w.Flush()
}
