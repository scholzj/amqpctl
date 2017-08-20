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
	"strings"
)

var readType string

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read attributeName/attributeValue",
	Example: "amqpctl read name/myListener1",
	Short: "Retrieve the attributes of a Manageable Entity.",
	Long: `Retrieve the attributes of a Manageable Entity. The entity is specified using an argument in the form of attributeName/attributeValue (e.g. name/myListener).`,
	Run: func(cmd *cobra.Command, args []string) {
		read(args)
	},
}

func init() {
	RootCmd.AddCommand(readCmd)
	readCmd.Flags().StringVar(&readType,"type","", "Type of the Manageable entity which should be read")
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

	var output bytes.Buffer
	output, err = read_operation.Read(&link, readType, attributeName, attributeValue)

	if err == nil {
		fmt.Print(output.String())
	} else {
		fmt.Printf("AMQP Management operation failed: %v\n", err.Error())
		os.Exit(1)
	}
}