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
	"github.com/labring-actions/gh-rebot/pkg/template"
	"github.com/labring-actions/gh-rebot/pkg/types"
	"strings"
)

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

func Changelog(reviews []string) error {
	branchName := generateBranchName()
	branchName = fmt.Sprintf("changelog-%s", branchName)
	if err := setPreGithub(); err != nil {
		return err
	}

	shells := []any{
		fmt.Sprintf(newBranch, branchName),
		types.RetryShell(fmt.Sprintf(generateChangelog, template.TryParseString(types.GlobalsBotConfig.GetChangelogScript(), types.GlobalsBotConfig))),
	}
	if err := types.ExecShellForAny()(shells); err != nil {
		return err
	}
	if release, ok, err := checkAndCommit(); err != nil {
		return err
	} else {
		if ok {
			copilot := ""
			if types.GlobalsBotConfig.Bot.Copilot4prs {
				copilot = "<br/>copilot:all"
			}
			afterShell := []any{
				fmt.Sprintf(gitPush, branchName),
				template.TryParseString(gitPRTmpl, map[string]string{
					"Title":     "docs: Automated Changelog Update for " + release,
					"Body":      "🤖 add release changelog using rebot." + copilot,
					"Reviewers": strings.Join(reviews, ","),
				}),
			}
			if err = types.ExecShellForAny()(afterShell); err != nil {
				return err
			}
		}
	}
	return nil
}
