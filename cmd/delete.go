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
	delete_operation "github.com/scholzj/amqpctl/operation/delete"
	"strings"
)

var deleteType string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete attributeName/attributeValue",
	Example: "amqpctl delete name/myListener1",
	Short: "Delete a Manageable Entity.",
	Long: `Delete a Manageable Entity. The entity is specified using an argument in the form of attributeName/attributeValue (e.g. name/myListener).`,
	Run: func(cmd *cobra.Command, args []string) {
		delete(args)
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&deleteType,"type","", "Type of the Manageable entity which should be deleted")
}

func delete(args []string) {
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

	var attributeName string
	var attributeValue string
	if len(args) > 0 && strings.Contains(args[0], "/"){
		readPair := strings.SplitN(args[0], "/", 2)

		if readPair[0] == "identity" {
			attributeName = "identity"
			attributeValue = readPair[1]
		} else if readPair[0] == "name" {
			attributeName = "name"
			attributeValue = readPair[1]
		} else {
			// WD 11 allows to query all attributes
			attributeName = readPair[0]
			attributeValue = readPair[1]
		}
	} else {
		fmt.Printf("Missing argument: The entity has to be specified using an argument in the form of attributeName/attributeValue (e.g. name/myListener)!\n")
		os.Exit(1)
	}

	err = delete_operation.Delete(&link, readType, attributeName, attributeValue)

	if err == nil {
		fmt.Print("Manageable Entity successfully deleted.\n")
	} else {
		fmt.Printf("AMQP Management operation failed: %v\n", err.Error())
		os.Exit(1)
	}
}