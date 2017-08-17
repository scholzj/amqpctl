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
)

var deleteAttributeName string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete identity/key",
	Short: "Delete a Manageable Entity.",
	Long: `Delete a Manageable Entity.`,
	Run: func(cmd *cobra.Command, args []string) {
		delete(args)
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVar(&deleteAttributeName,"attribute","identity", "Delete based on specific attribute (index)")
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

	var identityOrKey string
	if len(args) > 0 {
		identityOrKey = args[0]
	} else {
		fmt.Printf("Identity must be specified!\n")
		os.Exit(1)
	}

	err = delete_operation.Delete(&link, identityOrKey, deleteAttributeName)

	if err == nil {
		fmt.Print("Manageable Entity successfully deleted.\n")
	} else {
		fmt.Printf("AMQP Management operation failed: %v\n", err.Error())
		os.Exit(1)
	}
}