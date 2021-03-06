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
	"os"
	"github.com/scholzj/amqpctl/mgmtlink"
	"bytes"
	get_operation "github.com/scholzj/amqpctl/operation/get"
)

// typesCmd represents the types command
var typesCmd = &cobra.Command{
	Use:   "types [entityType]",
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

	var output bytes.Buffer

	if len(args) > 0 {
		output, err = get_operation.GetTypes(&link, args[0])
	} else {
		output, err = get_operation.GetTypes(&link, "")
	}

	if err == nil {
		fmt.Print(output.String())
	} else {
		fmt.Printf("AMQP Management operation failed: %v\n", err.Error())
		os.Exit(1)
	}
}