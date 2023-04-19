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
	"github.com/labring-actions/gh-rebot/pkg/config"
	"github.com/labring-actions/gh-rebot/pkg/template"
	"github.com/labring-actions/gh-rebot/pkg/utils"
	"k8s.io/client-go/util/retry"
	"strings"
)

type RetryShell string
type RetrySecretShell string
type SecretShell string

var execFn = func(shells []any) error {
	for _, sh := range shells {
		if s, ok := sh.(RetryShell); ok {
			return retry.RetryOnConflict(retry.DefaultRetry, func() error {
				return utils.RunCommand("bash", "-c", string(s))
			})
		}
		if s, ok := sh.(RetrySecretShell); ok {
			return retry.RetryOnConflict(retry.DefaultRetry, func() error {
				return utils.RunCommandInSecret(string(s), config.GlobalsConfig.GetToken())
			})
		}
		if s, ok := sh.(SecretShell); ok {
			if err := utils.RunCommandInSecret(string(s), config.GlobalsConfig.GetToken()); err != nil {
				return err
			}
		}
		if s, ok := sh.(string); ok {
			if err := utils.RunCommand("bash", "-c", s); err != nil {
				return err
			}
		}
	}
	return nil
}

func setPreGithub() error {
	shells := []any{
		authStatus,
		disablePrompt,
		fmt.Sprintf(forkRepo, config.GlobalsConfig.GetRepoName(), config.GlobalsConfig.GetForkName(), config.GlobalsConfig.GetOrgCommand()),
		RetryShell(fmt.Sprintf(checkRepo, config.GlobalsConfig.GetRepoName())),
		RetryShell(fmt.Sprintf(cloneRepo, config.GlobalsConfig.GetRepoName())),
		fmt.Sprintf(configEmail, config.GlobalsConfig.GetEmail()),
		fmt.Sprintf(configUser, config.GlobalsConfig.GetUsername()),
		SecretShell(fmt.Sprintf(setToken, config.GlobalsConfig.GetUsername(), config.GlobalsConfig.GetToken(), config.GlobalsConfig.GetRepoName())),
		SecretShell(fmt.Sprintf(gitAddRemote, config.GlobalsConfig.GetUsername(), config.GlobalsConfig.GetToken(), config.GlobalsConfig.GetForkName())),
		fmt.Sprintf(syncRepo),
	}
	if err := execFn(shells); err != nil {
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
		fmt.Sprintf(generateChangelog, template.TryParseString(config.GlobalsConfig.GetChangelogScript(), config.GlobalsConfig)),
	}
	if err := execFn(shells); err != nil {
		return err
	}
	if release, ok, err := checkAndCommit(); err != nil {
		return err
	} else {
		if ok {
			afterShell := []any{fmt.Sprintf(gitPush, branchName), fmt.Sprintf(gitPR, release, strings.Join(reviews, ","))}
			if err = execFn(afterShell); err != nil {
				return err
			}
		}
	}
	return nil
}
