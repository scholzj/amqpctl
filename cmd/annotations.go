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
	get_operation "github.com/scholzj/amqpctl/operation/get"
)

// annotationsCmd represents the annotations command
var annotationsCmd = &cobra.Command{
	Use:   "annotations [entityType]",
	Aliases: []string{"annotation"},
	Short: "Retrieve the list of Type Annotations",
	Long: `Retrieve the list of Type Annotations implemented by each Manageable Entity Type.`,
	Run: func(cmd *cobra.Command, args []string) {
		getAnnotations(args)
	},
}

func init() {
	getCmd.AddCommand(annotationsCmd)
}

func getAnnotations(args []string) {
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
		output, err = get_operation.GetAnnotations(&link, args[0])
	} else {
		output, err = get_operation.GetAnnotations(&link, "")
	}

	if err == nil {
		fmt.Print(output.String())
	} else {
		fmt.Printf("AMQP Management operation failed: %v\n", err.Error())
		os.Exit(1)
	}
}
