// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"text/tabwriter"
	"os"
	"qpid.apache.org/amqp"
	"strings"
	"github.com/scholzj/amqpctl/mgmtlink"
)

// typesCmd represents the types command
var typesCmd = &cobra.Command{
	Use:   "types",
	Aliases: []string{"type"},
	Short: "Get list of Manageable Entity Types",
	Long: `Get list of Manageable Entity Types that can be managed via this Management Node`,
	Run: func(cmd *cobra.Command, args []string) {
		getTypes(args)
	},
}

func init() {
	getCmd.AddCommand(typesCmd)
}

func getTypes(args []string) {
	link := mgmtlink.AmqpMgmtLink{}
	err := link.ConfigureConnection(amqpCfg)
	if err != nil {
		fmt.Printf("Failed to configure AMQP connection: %v\n", err.Error())
		os.Exit(1)
	}

	err = link.Connect()
	if err != nil {
		fmt.Printf("Failed to connect to AMQP endpoint: %v\n", err.Error())
		os.Exit(1)
	}

	defer link.Close()

	var reqProperties map[string]interface{}

	if len(args) > 0 {
		reqProperties = map[string]interface{}{"operation": "GET-TYPES", "entityType": args[0]}
	} else {
		reqProperties = map[string]interface{}{"operation": "GET-TYPES"}
	}

	respProperties, respBody, err := link.Operation(reqProperties, nil)

	if err == nil {
		if respProperties["statusCode"].(int64) == 200 {
			printTypes(respProperties, respBody)
		} else {
			fmt.Printf("AMQP Management operation wsn't successfull: %v (%v)\n", respProperties["statusCode"], respProperties["statusDescription"])
			os.Exit(1)
		}
	} else {
		fmt.Printf("AMQP Management operation failed: %v\n", err.Error())
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