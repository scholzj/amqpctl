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
	create_operation "github.com/scholzj/amqpctl/operation/create"
	"strings"
)

// readCmd represents the read command
var createCmd = &cobra.Command{
	Use:   "create attribute1=value [attribute2=value...]",
	Short: "Create a new Manageable Entity.",
	Long: `Create a new Manageable Entity.`,
	Run: func(cmd *cobra.Command, args []string) {
		create(args)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}

func create(args []string) {
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

	createMap := make(map[interface{}]interface{})
	if len(args) > 0 {
		for _, change := range args {
			createPair := strings.SplitN(change, "=", 2)
			createMap[createPair[0]] = createPair[1]
		}
	} else {
		fmt.Printf("At least one attribue=value pair has to be specified!\n")
		os.Exit(1)
	}

	var output bytes.Buffer
	output, err = create_operation.Create(&link, createMap)

	if err == nil {
		fmt.Print(output.String())
	} else {
		fmt.Printf("AMQP Management operation failed: %v\n", err.Error())
		os.Exit(1)
	}
}