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
)

// typesCmd represents the types command
var typesCmd = &cobra.Command{
	Use:   "types",
	Aliases: []string{"type"},
	Short: "Get list of Manageable Entity Types",
	Long: `Get list of Manageable Entity Types that can be managed via this Management Node`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("types called")
	},
}

func init() {
	getCmd.AddCommand(typesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// typesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// typesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
