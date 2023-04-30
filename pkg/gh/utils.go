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

package gh

import (
	"fmt"
	"github.com/cuisongliu/logger"
	"github.com/labring/gh-rebot/pkg/types"
	"github.com/labring/gh-rebot/pkg/utils"
	"math/rand"
	"strings"
	"time"
)

func generateBranchName() string {
	timestamp := time.Now().Format("20060102")
	randomCode := randString(6)
	return fmt.Sprintf("%s-%s", timestamp, randomCode)
}

func randString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func checkAndCommit() (string, bool, error) {
	// Check workspace status
	statusOutput, err := utils.RunCommandWithOutput(gitStatus, true)
	if err != nil {
		return "", false, fmt.Errorf("failed to get workspace status: %v", err)
	}
	release, err := utils.RunCommandWithOutput(gitRelease, true)
	if err != nil {
		return "", false, err
	}
	// Determine if there are changes to be committed
	if strings.Contains(statusOutput, "Changes not staged for commit") || strings.Contains(statusOutput, "Changes to be committed") {
		// Add all changes to staging area
		err = utils.RunCommand("bash", "-c", gitAdd)
		if err != nil {
			return "", false, fmt.Errorf("failed to add changes to staging area: %v", err)
		}
		commitMessage := fmt.Sprintf("Automated Changelog Update: Update directory for %s release", release)
		err = utils.RunCommand("bash", "-c", fmt.Sprintf(gitCommit, commitMessage))
		if err != nil {
			return "", false, fmt.Errorf("commit failed: %v", err)
		}
		fmt.Println("Changes committed")
		return release, true, nil
	}
	fmt.Println("No changes to be committed")
	return "", false, nil
}

func checkRemoteTagExists(tag string) (bool, error) {
	statusOutput, err := utils.RunCommandWithOutput(gitTag, false)
	if err != nil {
		return false, fmt.Errorf("failed to list tags: %v", err)
	}
	tags := strings.Split(statusOutput, "\n")
	for _, t := range tags {
		if t == tag {
			return true, nil
		}
	}

	return false, nil
}

func setPreGithub() error {
	shells := []any{
		authStatus,
		disablePrompt,
		fmt.Sprintf(forkRepo, types.GlobalsBotConfig.GetRepoName(), types.GlobalsBotConfig.GetForkName(), types.GlobalsBotConfig.GetOrgCommand()),
		types.RetryShell(fmt.Sprintf(checkRepo, types.GlobalsBotConfig.GetRepoName())),
		types.RetryShell(fmt.Sprintf(cloneRepo, types.GlobalsBotConfig.GetRepoName())),
		fmt.Sprintf(configEmail, types.GlobalsBotConfig.GetEmail()),
		fmt.Sprintf(configUser, types.GlobalsBotConfig.GetUsername()),
		types.SecretShell(fmt.Sprintf(setToken, types.GlobalsBotConfig.GetUsername(), types.GlobalsBotConfig.GetToken(), types.GlobalsBotConfig.GetRepoName())),
		types.SecretShell(fmt.Sprintf(gitAddRemote, types.GlobalsBotConfig.GetUsername(), types.GlobalsBotConfig.GetToken(), types.GlobalsBotConfig.GetForkName())),
		fmt.Sprintf(syncRepo),
	}
	if err := types.ExecShellForAny(types.GlobalsBotConfig.GetToken())(shells); err != nil {
		logger.Error("setPreGithub err:%v", err)
		return err
	}
	return nil
}
