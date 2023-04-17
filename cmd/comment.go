/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/cuisongliu/logger"
	"github.com/labring-actions/gh-rebot/pkg/config"
	"github.com/labring-actions/gh-rebot/pkg/gh"
	"github.com/labring-actions/gh-rebot/pkg/utils"
	"strings"

	"github.com/spf13/cobra"
)

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:  "comment",
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		safeRepo := gh.GlobalsGithubVar.SafeRepo
		runnerID := gh.GlobalsGithubVar.RunnerID
		comment := gh.GlobalsGithubVar.CommentBody
		issueID := gh.GlobalsGithubVar.IssueOrPRNumber
		runnerURL := fmt.Sprintf("https://github.com/%s/actions/runs/%s", safeRepo, runnerID)
		switch {
		case strings.HasPrefix(comment, "/sealos_bot_release"):
			data := strings.Split(comment, " ")
			if len(data) == 2 && utils.ValidateVersion(data[1]) {

				msg := config.GlobalsConfig.GetMessage("release_success", "release action finished successfully!")
				return gh.SendMsgToIssue(issueID, msg, runnerURL, safeRepo)
			} else {
				msg := config.GlobalsConfig.GetMessage("release_format_error", "release action failed!")
				logger.Error("command format is error: %s ex. /sealos_bot_release {tag}", comment)
				return gh.SendMsgToIssue(issueID, msg, runnerURL, safeRepo)
			}
		}
		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		printEnvs()
		var err error
		gh.GlobalsGithubVar, err = gh.GetGHEnvToVar()
		if err != nil {
			return err
		}
		logger.Debug("github env to var: %v", gh.GlobalsGithubVar)
		if err = checkPermission(); err != nil {
			return err
		}
		if err = checkGithubEnv(); err != nil {
			return err
		}
		return nil
	},
}

func checkGithubEnv() error {
	if gh.GlobalsGithubVar.RunnerID == "" {
		return fmt.Errorf("error: GITHUB_RUN_ID is not set. Please set the GITHUB_RUN_ID environment variable")
	}
	if gh.GlobalsGithubVar.SafeRepo == "" {
		return fmt.Errorf("error: not found repository.full_name in github event")
	}
	if gh.GlobalsGithubVar.CommentBody == "" {
		return fmt.Errorf("error: not found comment.body in github event")
	}
	if gh.GlobalsGithubVar.IssueOrPRNumber == "" {
		return fmt.Errorf("error: not found issue.number or pull_request.number in github event")
	}
	return nil
}

func init() {
	rootCmd.AddCommand(commentCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
