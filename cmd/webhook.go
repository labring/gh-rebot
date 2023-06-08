/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/cuisongliu/logger"
	"github.com/labring/gh-rebot/pkg/setup"
	"github.com/labring/gh-rebot/pkg/webhook"

	"github.com/spf13/cobra"
)

// webhookCmd represents the webhook command
var webhookCmd = &cobra.Command{
	Use: "webhook",
	Run: func(cmd *cobra.Command, args []string) {
		err := webhook.RegistryHttpServer(8080)
		if err != nil {
			logger.Fatal(err, "unable to init http server")
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		setup.Setup(cfgFile)
	},
}

func init() {
	rootCmd.AddCommand(webhookCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// webhookCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// webhookCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
