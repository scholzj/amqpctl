package main

import (
	"github.com/docopt/docopt-go"
	"fmt"
	"os"
	"./operations"
	"./utils"
	"io/ioutil"
	"crypto/x509"
	"crypto/tls"
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
  --hostname HOSTNAME		AMQP hostname (default localhost)
  --port PORT			AMQP port (default 5672)
  --username USERNAME   	AMQP username
  --password PASSWORD   	AMQP password
  --sasl-mechanism MECH 	AMQP SASL mechanism
  --ssl-ca CAFILE		SSL certification authority certificate(s)
  --ssl-skip-verify		Skip hostname verification
  --ssl-cert CERTFILE		SSL certificate for client authentication
  --ssl-key KEYFILE		SSL private key for client authentication
  -h --help             	Show this screen.

Description:
  The amqpctl command line tool implements AMQP Management specification

  See 'amqpctl <operation> --help' to read about a specific operations.
`
	arguments, _ := docopt.Parse(usage, nil, true, "0.0.1", false, false)

	if arguments["<operation>"] != nil {
		operation := arguments["<operation>"].(string)
		args := append([]string{operation}, arguments["<args>"].([]string)...)
		mgmtLink := parseConnectionArgs(arguments)

		switch operation {
		case "version":
			operations.Version(args)
		case "query":
			operations.Query(args, mgmtLink)
		case "gettypes":
			operations.GetTypes(args, mgmtLink)
		case "getattributes":
			operations.GetAttributes(args, mgmtLink)
		case "getoperations":
			operations.GetOperations(args, mgmtLink)
		case "getannotations":
			operations.GetAnnotations(args, mgmtLink)
		case "getmgmtnodes":
			operations.GetMgmtNodes(args, mgmtLink)
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

func parseConnectionArgs(args map[string]interface{}) (mgmtLink utils.MgmtLink) {
	mgmtLink = utils.MgmtLink{}

	// URL (Hostname / port)
	hostname := args["--hostname"]
	if hostname == nil {
		hostname = "localhost"
	}

	port := args["--port"]
	if port == nil {
		port = "5672"
	}

	mgmtLink.Url = fmt.Sprintf("%v:%v", hostname, port)

	// Username
	if username := args["--username"]; username != nil {
		mgmtLink.Username = username.(string)
	}

	// Password
	if password := args["--password"]; password != nil {
		mgmtLink.Password = password.(string)
	}

	// SASL Mechanism
	if sasl := args["--sasl-mechanism"]; sasl != nil {
		mgmtLink.SaslMechanism = sasl.(string)
	}

	// CA certificate
	if ca := args["--ssl-ca"]; ca != nil {
		brokerCert, err := ioutil.ReadFile(ca.(string))
		if err != nil {
			fmt.Printf("Ups, something went wrong while loading CA certificate ... %s", err)
			os.Exit(1)
		}

		brokerCertPool := x509.NewCertPool()
		brokerCertPool.AppendCertsFromPEM(brokerCert)

		mgmtLink.BrokerCertificate = brokerCertPool

		// Hostname verification
		if verify := args["--ssl-skip-verify"]; verify != nil {
			mgmtLink.SslSkipVerify = true
		} else {
			mgmtLink.SslSkipVerify = false
		}

		// Client certificate and key
		cert := args["--ssl-cert"]
		key := args["--ssl-key"];
		if cert != nil && key != nil {
			memberKey, err := tls.LoadX509KeyPair(cert.(string), key.(string))
			if err != nil {
				fmt.Printf("Ups, something went wrong while loading client certificate / key ... %s", err)
				os.Exit(1)
			}

			mgmtLink.ClientCertificate = &memberKey
		}
	}





	return
}