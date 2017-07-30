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
	"github.com/scholzj/amqpctl/mgmtlink"
	"os"
	"bytes"
	"github.com/scholzj/amqpctl/operation"
)

var allAttributes bool

// queryCmd represents the query command
var queryCmd = &cobra.Command{
	Use:   "query [entityType] [attributes...]",
	Short: "Retrieve selected attributes of Manageable Entities",
	Long: `Retrieve selected attributes of Manageable Entities that can be read at this Management Node.`,
	Run: func(cmd *cobra.Command, args []string) {
		query(args)
	},
}

func init() {
	RootCmd.AddCommand(queryCmd)
	queryCmd.Flags().BoolVar(&allAttributes,"all-attributes",false, "Query all available attributes")
}

func query(args []string) {
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

	var entityName string
	if len(args) > 0 {
		entityName = args[0]
	} else {
		entityName = ""
	}

	var attributes []string
	if allAttributes {
		attributes = []string{}
	} else if len(args) > 1 {
		attributes = args[1:]
	} else {
		attributes = []string{"name", "type", "identity"}
	}

	var output bytes.Buffer
	output, err = operation.Query(&link, entityName, attributes)

	if err == nil {
		fmt.Print(output.String())
	} else {
		fmt.Printf("AMQP Management operation failed: %v\n", err.Error())
		os.Exit(1)
	}
}