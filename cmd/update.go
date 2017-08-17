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
	update_operation "github.com/scholzj/amqpctl/operation/update"
	"strings"
)

var updateAttributeName string

// readCmd represents the read command
var updateCmd = &cobra.Command{
	Use:   "update identity/key attribute1=value [attribute2=value...]",
	Short: "Update a Manageable Entity.",
	Long: `Update a Manageable Entity.`,
	Run: func(cmd *cobra.Command, args []string) {
		update(args)
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringVar(&updateAttributeName,"attribute","identity", "Update entity based on specific attribute (index)")
}

func update(args []string) {
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

	var identityOrKey string
	var changeMap map[interface{}]interface{}
	changeMap = make(map[interface{}]interface{})
	if len(args) > 0 {
		identityOrKey = args[0]

		if len(args) > 1 {
			changes := args[1:]

			for _, change := range changes {
				changePair := strings.SplitN(change, "=", 2)
				changeMap[changePair[0]] = changePair[1]
			}

		} else {
			fmt.Printf("At least one attribue=value pair has to be specified!\n")
			os.Exit(1)
		}

	} else {
		fmt.Printf("Identity must be specified!\n")
		os.Exit(1)
	}

	var output bytes.Buffer
	output, err = update_operation.Update(&link, identityOrKey, updateAttributeName, changeMap)

	if err == nil {
		fmt.Print(output.String())
	} else {
		fmt.Printf("AMQP Management operation failed: %v\n", err.Error())
		os.Exit(1)
	}
}