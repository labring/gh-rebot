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
	"bytes"
	"fmt"
	"github.com/labring-actions/gh-rebot/pkg/template"
	"github.com/labring-actions/gh-rebot/pkg/types"
	"github.com/labring-actions/gh-rebot/pkg/utils"
	"github.com/pkg/errors"
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

func SendMsgToIssue(msg string, actionURL ...string) error {
	tpl, ok, _ := template.TryParse(`gh issue comment {{.IssueOrPRNumber}} --body "{{.Msg}} <br/>See: <br/>- {{.GetRunnerURL}}{{range .ActionURLs}}<br/>- {{.}}{{end}}" --repo {{.SafeRepo}}`)
	if ok {
		out := bytes.NewBuffer(nil)
		_ = tpl.Execute(out, map[string]interface{}{
			"IssueOrPRNumber": types.GlobalsGithubVar.IssueOrPRNumber,
			"Msg":             msg,
			"GetRunnerURL":    types.GlobalsGithubVar.GetRunnerURL(),
			"SafeRepo":        types.GlobalsGithubVar.SafeRepo,
			"ActionURLs":      actionURL,
		})
		return utils.RunCommand("bash", "-c", out.String())
	}

	return errors.New("template parse error")
}

func SendCustomizeMsgToIssue(msg string) error {
	tpl, ok, _ := template.TryParse(`gh issue comment {{.IssueOrPRNumber}} --body "{{.Msg}}" --repo {{.SafeRepo}}`)
	if ok {
		out := bytes.NewBuffer(nil)
		_ = tpl.Execute(out, map[string]interface{}{
			"IssueOrPRNumber": types.GlobalsGithubVar.IssueOrPRNumber,
			"Msg":             msg,
			"SafeRepo":        types.GlobalsGithubVar.SafeRepo,
		})
		return utils.RunCommand("bash", "-c", out.String())
	}

	return errors.New("template parse error")
}
