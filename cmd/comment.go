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
	"github.com/cuisongliu/logger"
	"github.com/labring-actions/gh-rebot/pkg/types"
	"github.com/labring-actions/gh-rebot/pkg/workflow"
	"strings"

	"github.com/spf13/cobra"
)

var commentCmd = &cobra.Command{
	Use:  "comment",
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		comment := types.GlobalsGithubVar.CommentBody
		cmds := strings.Split(comment, "\n")
		for _, t := range cmds {
			logger.Debug("cmds: ", strings.TrimSpace(t))
			wfs := make([]workflow.Interface, 0)
			if types.GlobalsBotConfig.Release != nil {
				wfs = append(wfs, workflow.NewRelease(strings.TrimSpace(t)))
			}
			used := 0
			for _, wf := range wfs {
				if strings.HasPrefix(t, wf.Comment()) {
					if err := wf.Run(); err != nil {
						return err
					}
					used++
				}
			}
			if used == 0 {
				logger.Warn("not support command: ", t)
			}
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := preCheck(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	cobra.OnInitialize(func() {
		checkToken()
	})
	rootCmd.AddCommand(commentCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
