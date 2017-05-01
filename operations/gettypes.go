package operations

import (
	"github.com/docopt/docopt-go"
	"fmt"
	"strings"
	"os"
	"qpid.apache.org/electron"
	"qpid.apache.org/amqp"
	"time"
	"text/tabwriter"
)

func GetTypes(args []string) {
	usage := `Usage:
  amqpctl gettypes [<entityType>]

Options:
  -h --help   Show this screen.

Description:
  List the types supported by the management node.
`
	arguments, err := docopt.Parse(usage, args, true, "", false, false)
	if err != nil {
		fmt.Printf("Invalid option: 'amqpctl %s'. Use flag '--help' to read about a specific subcommand.\n", strings.Join(args, " "))
		os.Exit(1)
	}
	if len(arguments) == 0 {
		return
	}

	// Get types from the management endpoint
	container := electron.NewContainer(fmt.Sprintf("myContainer"))

	conn, err := container.Dial("tcp", "localhost:5672")
	if err != nil {
		fmt.Printf("Failed to open AMQP connection: %s", err)
		os.Exit(1)
	}

	sess, err := conn.Session(electron.IncomingCapacity(1024000), electron.OutgoingWindow(1024))
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		os.Exit(1)
	}

	recv, err := sess.Receiver(electron.Source("tempAddress"), electron.Capacity(100), electron.Prefetch(true))
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		os.Exit(1)
	}

	snd, err := sess.Sender(electron.Target("$management"))
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		os.Exit(1)
	}

	m := amqp.NewMessage()
	m.SetReplyTo("tempAddress")
	if arguments["<entityType>"] != nil {
		m.SetProperties(map[string]interface{}{"operation": "GET-TYPES", "entityType": arguments["<entityType>"]})
	} else {
		m.SetProperties(map[string]interface{}{"operation": "GET-TYPES"})
	}
	//
	body := map[string]interface{}{"attributeNames": []string{}}
	m.Marshal(body)

	res := snd.SendSync(m)
	if res.Error != nil {
		fmt.Printf("Failed to send message: %s\n", res.Error)
		os.Exit(1)
	}

	msg, err := recv.ReceiveTimeout(time.Duration(10 * time.Second))
	//msg, err := recv.Receive()

	if err == nil {
		//fmt.Printf("Message received: %s\n", msg.Message.Body())
		parseResult(msg.Message)
		msg.Accept()

	} else if err == electron.Timeout {
		fmt.Printf("No message received\n")
	} else {
		fmt.Printf("Something went wrong: %s\n", err)
		os.Exit(1)
	}

	conn.Close(nil)
}

func parseResult(m amqp.Message) {
	// Parse the properties
	/*props := m.Properties()
	fmt.Printf("Status code:\t%v\n", props["statusCode"])*/
	/*for k, v := range props {
		fmt.Printf("    %v: %v\n", k, v)
	}*/

	types := m.Body().(amqp.Map)
	/*types := body["attributeNames"].(amqp.List)
	res := body["results"].(amqp.List)

	for _, listener := range res {
		fmt.Print("Connectors:\n")
		for i, v := range listener.(amqp.List) {
			//if v != nil {
			fmt.Printf("    %s: %v\n", atts[i], v)
			//}
		}
	}*/

	w := tabwriter.NewWriter(os.Stdout, 10, 4, 3, ' ', 0)
	fmt.Fprint(w, "TYPE\tPARENTS\t\n")

	for entitytype, extends := range types {
		parents := make([]string, len(extends.(amqp.List)))
		for i, parent := range extends.(amqp.List) {
			parents[i] = parent.(string)
		}


		fmt.Fprintf(w, "%v\t%v\t\n", entitytype, strings.Join(parents, ", "))
	}

	w.Flush()
}

func printResult() {

}
