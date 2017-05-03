package main

import (
	"github.com/docopt/docopt-go"
	"fmt"
	"os"
	"./operations"
)

func main() {
	usage := `Usage:
  amqpctl [options] <operation> [<args>...]

    query 	   	Query selected attributes of Management entities
    gettypes    	Get list of supported types
    getattributes   	Get list of attributes supported by entity types
    getoperations   	Get list of operations supported by entity types
    getannotations   	Get list of annotations supported by entity types
    getmgmtnodes	Get list of other management nodes
    version   		Display the version of amqpctl.

Options:
  -h --help               Show this screen.
  -l --log-level=<level>  Set the log level (one of panic, fatal, error,
                          warn, info, debug) [default: panic]

Description:
  The amqpctl command line tool implements AMQP Management specification

  See 'amqpctl <operation> --help' to read about a specific operations.
`
	arguments, _ := docopt.Parse(usage, nil, true, "0.0.1", true, false)

	if arguments["<operation>"] != nil {
		operation := arguments["<operation>"].(string)
		args := append([]string{operation}, arguments["<args>"].([]string)...)

		switch operation {
		case "version":
			operations.Version(args)
		case "query":
			operations.Query(args)
		case "gettypes":
			operations.GetTypes(args)
		case "getattributes":
			operations.GetAttributes(args)
		case "getoperations":
			operations.GetOperations(args)
		case "getannotations":
			operations.GetAnnotations(args)
		case "getmgmtnodes":
			operations.GetMgmtNodes(args)
		/*case "create":
			commands.Create(args)
		case "replace":
			commands.Replace(args)
		case "apply":
			commands.Apply(args)
		case "delete":
			commands.Delete(args)
		case "get":
			commands.Get(args)
		case "version":
			commands.Version(args)
		case "node":
			commands.Node(args)
		case "ipam":
			commands.IPAM(args)
		case "config":
			commands.Config(args)*/
		default:
			fmt.Printf("Unknown operation: %q\n", operation)
			fmt.Println(usage)
			os.Exit(1)
		}
	}
}