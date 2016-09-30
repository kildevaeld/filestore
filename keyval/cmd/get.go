// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
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
	"io"
	"os"

	"github.com/kildevaeld/filestore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get <key>",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here

		var (
			store  filestore.Store
			reader io.ReadCloser
			err    error
		)

		if store, err = getFileStore(); err != nil {
			exit(1, err)
		}

		output := os.Stdout
		outstr := viper.GetString("output")
		if outstr != "" {
			if output, err = os.Create(outstr); err != nil {
				exit(1, err)
			}
			defer output.Close()
		}

		if reader, err = store.Get([]byte(args[0])); err != nil {
			exit(1, err)
		}
		defer reader.Close()
		if _, err = io.Copy(output, reader); err != nil {
			exit(1, err)
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	getCmd.Flags().StringP("output", "o", "", "Output")

	flags := getCmd.Flags()
	viper.BindPFlag("output", flags.Lookup("output"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
