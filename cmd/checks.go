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
	"github.com/cuisongliu/logger"
	"github.com/labring-actions/gh-rebot/pkg/types"
	"github.com/pkg/errors"
	"os"
)

func preCheck() error {
	if err := checkGithubEnv(); err != nil {
		return err
	}
	return nil
}

func checkGithubEnv() error {
	if types.GlobalsGithubVar.RunnerID == "" {
		return fmt.Errorf("error: GITHUB_RUN_ID is not set. Please set the GITHUB_RUN_ID environment variable")
	}
	if types.GlobalsGithubVar.SafeRepo == "" {
		return fmt.Errorf("error: not found repository.full_name in github event")
	}
	if types.GlobalsGithubVar.CommentBody == "" {
		return fmt.Errorf("error: not found comment.body in github event")
	}
	if types.GlobalsGithubVar.IssueOrPRNumber == 0 {
		return fmt.Errorf("error: not found issue.number or pull_request.number in github event")
	}
	return nil
}

func checkToken() {
	var err error
	types.GlobalsGithubVar, err = types.GetGHEnvToVar()
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	logger.Debug("github env to var: %v", types.GlobalsGithubVar)
	types.GlobalsBotConfig, err = types.LoadConfig(cfgFile)
	if err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	if err = types.GlobalsBotConfig.Validate(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
	logger.Cfg(types.GlobalsBotConfig.GetDebug(), false)
	if err := checkGhToken(); err != nil {
		logger.Error(err)
		os.Exit(1)
	}
}

func checkGhToken() error {
	if _, ok := os.LookupEnv("GH_TOKEN"); !ok {
		return errors.New("error: GH_TOKEN is not set. Please set the GH_TOKEN environment variable to enable authentication and access to the GitHub API")
	}
	return nil
}
