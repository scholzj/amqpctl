package operations

import (
	"github.com/docopt/docopt-go"
	"fmt"
	"strings"
	"os"
)

func Version(args []string) {
	usage := `Usage:
  amqpctl version

Options:
  -h --help   Show this screen.

Description:
  Display the version of amqpctl.
`
	arguments, err := docopt.Parse(usage, args, true, "", false, false)
	if err != nil {
		fmt.Printf("Invalid option: 'amqpctl %s'. Use flag '--help' to read about a specific subcommand.\n", strings.Join(args, " "))
		os.Exit(1)
	}
	if len(arguments) == 0 {
		return
	}

	fmt.Println("0.0.1")
}
