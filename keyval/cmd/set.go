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
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/kildevaeld/filestore"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set <key> [value]",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			store  filestore.Store
			reader io.ReadCloser
			err    error
		)

		if store, err = getFileStore(); err != nil {
			exit(2, err)
		}

		file := viper.GetString("file")

		if file != "" {
			reader, err = os.Open(file)
			if err != nil {
				exit(2, err)
			}

		} else if len(args) > 1 {
			reader = ioutil.NopCloser(bytes.NewReader([]byte(args[1])))
		} else {
			exit(1, errors.New("You have to use --file or supply a value"))
		}
		defer reader.Close()
		if err = store.Set([]byte(args[0]), reader); err != nil {
			exit(1, err)
		}

	},
}

func init() {
	RootCmd.AddCommand(setCmd)

	flags := setCmd.Flags()
	flags.StringP("file", "f", "", "Input file")
	viper.BindPFlag("file", flags.Lookup("file"))

}
