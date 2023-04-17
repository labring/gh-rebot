/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/labring-actions/gh-rebot/pkg/gh"
	"github.com/labring-actions/gh-rebot/pkg/workflow"
	"strings"

	"github.com/spf13/cobra"
)

// commentCmd represents the comment command
var commentCmd = &cobra.Command{
	Use:  "comment",
	Args: cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		comment := gh.GlobalsGithubVar.CommentBody
		wf := workflow.NewWorkflow(comment)
		switch {
		case strings.HasPrefix(comment, "/sealos_bot_release"):
			return wf.Release()
		case strings.HasPrefix(comment, "/sealos_bot_changelog"):
			return wf.Changelog()
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
	rootCmd.AddCommand(commentCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
