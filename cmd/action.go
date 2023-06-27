/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/cuisongliu/logger"
	"github.com/labring/gh-rebot/pkg/action"
	"github.com/labring/gh-rebot/pkg/setup"
	"github.com/labring/gh-rebot/pkg/types"
	"github.com/labring/gh-rebot/pkg/utils"
	"os"

	"github.com/spf13/cobra"
)

// actionCmd represents the action command
var actionCmd = &cobra.Command{
	Use: "action",
	Run: func(cmd *cobra.Command, args []string) {
		actionType, err := action.GetEnvFromAction("type")
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
		logger.Debug("action type is ", actionType)
		if !utils.In([]string{"issue_renew"}, actionType) {
			if err := types.ValidateIssueOrPRNumber(); err != nil {
				logger.Error(err)
				os.Exit(1)
			}
		}
		switch actionType {
		case "/comment":
			err = action.CommentEngine()
		case "issue_comment":
			err = action.IssueComment()
		case "issue_renew":
			err = action.IssueRenew()
		default:
			err = fmt.Errorf("not support action type")
		}
		if err != nil {
			logger.Error(err)
			os.Exit(1)
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		setup.Setup(cfgFile)
	},
}

func init() {
	rootCmd.AddCommand(actionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// actionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// actionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
