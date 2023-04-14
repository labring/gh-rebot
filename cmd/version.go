/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"github.com/labring-actions/gh-rebot/pkg/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/json"
	"sigs.k8s.io/yaml"
)

var shortPrint bool
var output string

func newVersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:     "version",
		Short:   "Print version info",
		Args:    cobra.NoArgs,
		Example: `runtime-ctl version`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// output default to be yaml
			if output != "yaml" && output != "json" {
				return errors.New(`--output must be 'yaml' or 'json'`)
			}

			if shortPrint {
				fmt.Println(version.Get().String())
				return nil
			}
			if err := PrintInfo(); err != nil {
				return err
			}
			return nil
		},
	}
	versionCmd.Flags().BoolVar(&shortPrint, "short", false, "if true, print just the version number.")
	versionCmd.Flags().StringVarP(&output, "output", "o", "yaml", "One of 'yaml' or 'json'")
	return versionCmd
}

func init() {
	rootCmd.AddCommand(newVersionCmd())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// versionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// versionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func PrintInfo() error {
	var (
		marshalled []byte
	)
	info := version.Get()
	var err error
	switch output {
	case "yaml":
		marshalled, err = yaml.Marshal(&info)
		if err != nil {
			return fmt.Errorf("fail to marshal yaml: %w", err)
		}
		fmt.Println(string(marshalled))
	case "json":
		marshalled, err = json.Marshal(&info)
		if err != nil {
			return fmt.Errorf("fail to marshal json: %w", err)
		}
		fmt.Println(string(marshalled))
	default:
		// There is a bug in the program if we hit this case.
		// However, we follow a policy of never panicking.
		return fmt.Errorf("versionOptions were not validated: --output=%q should have been rejected", output)
	}
	return nil
}
