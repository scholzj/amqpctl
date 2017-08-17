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
	read_operation "github.com/scholzj/amqpctl/operation/read"
)

var readAttributeName string

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read identity/key",
	Short: "Retrieve the attributes of a Manageable Entity.",
	Long: `Retrieve the attributes of a Manageable Entity.`,
	Run: func(cmd *cobra.Command, args []string) {
		read(args)
	},
}

func init() {
	RootCmd.AddCommand(readCmd)
	readCmd.Flags().StringVar(&readAttributeName,"attribute","identity", "Read based on specific attribute (index)")
}

func read(args []string) {
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

	var output bytes.Buffer
	output, err = read_operation.Read(&link, identityOrKey, readAttributeName)

	if err == nil {
		fmt.Print(output.String())
	} else {
		fmt.Printf("AMQP Management operation failed: %v\n", err.Error())
		os.Exit(1)
	}
}