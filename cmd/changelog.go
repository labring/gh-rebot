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
	"github.com/labring-actions/gh-rebot/pkg/config"
	"github.com/labring-actions/gh-rebot/pkg/gh"
	"github.com/labring-actions/gh-rebot/pkg/utils"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var reviews []string

// changelogCmd represents the changelog command
var changelogCmd = &cobra.Command{
	Use: "changelog",
	RunE: func(cmd *cobra.Command, args []string) error {
		return gh.Changelog(reviews)
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		printEnvs()
		if len(reviews) == 0 {
			reviews = []string{"cuisongliu"}
		}

		if err := checkPermission(); err != nil {
			return err
		}
		return nil
	},
}

func printEnvs() {
	envs := os.Environ()
	for _, env := range envs {
		if strings.HasPrefix(env, "SEALOS_SYS_") {
			logger.Info("sealos system env: %s ", env)
		}
	}
}
func checkPermission() error {
	ops := config.GlobalsConfig.GetAllowOps()
	if len(ops) == 0 {
		return errors.New("Error: config bot.triggers is not set. Please set the bot.allowOps to config yaml.")
	}
	if tigger, ok := os.LookupEnv("SEALOS_SYS_TRIGGER"); !ok {
		return errors.New("Error: SEALOS_SYS_TRIGGER is not set. Please set the SEALOS_SYS_TRIGGER environment variable.")
	} else {
		if !utils.In(ops, tigger) {
			return errors.New("Error: no has permission to trigger this action.")
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(changelogCmd)
	changelogCmd.Flags().StringSliceVarP(&reviews, "reviews", "r", []string{}, "reviewers")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// changelogCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// changelogCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
