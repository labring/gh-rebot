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

package action

import (
	"context"
	"fmt"
	"github.com/cuisongliu/logger"
	"github.com/google/go-github/v39/github"
	github_go "github.com/labring/gh-rebot/pkg/github-go"
	"github.com/labring/gh-rebot/pkg/types"
	"os"
	"strings"
)

// PRComment is a action to comment on PR
func PRComment() error {
	fileName, err := GetEnvFromAction("filename")
	if err != nil {
		return err
	}
	replaceTag, err := GetEnvFromAction("replace_tag")
	if err != nil {
		return err
	}

	prNumber := int(types.ActionConfigJSON.IssueOrPRNumber)
	if err != nil {
		return err
	}

	fileContent, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}
	owner, repo, err := getRepo("")
	if err != nil {
		return err
	}
	logger.Debug("repo:%s, filename: %s, replaceTag: %s, prNumber: %d", repo, fileName, replaceTag, prNumber)
	ctx := context.Background()
	client := github_go.GithubClient(ctx)
	comments, _, err := client.Issues.ListComments(ctx, owner, repo, prNumber, nil)
	if err != nil {
		return fmt.Errorf("Issues.ListComments returned error: %v", err)
	}
	hiddenReplace := fmt.Sprintf("<!-- %s -->", replaceTag)
	content := string(fileContent) + "\n" + hiddenReplace

	// Checks existing comments, edits if match found
	for _, comment := range comments {
		if comment.Body != nil && comment.ID != nil {
			if *comment.Body == content {
				logger.Debug("The comment %d has been already added to the pull request. Skipping...", *comment.ID)
				return nil
			} else if hiddenReplace != "" && strings.LastIndex(*comment.Body, hiddenReplace) != -1 {
				_, _, err = client.Issues.EditComment(ctx, owner, repo, *comment.ID, &github.IssueComment{Body: github.String(content)})
				if err != nil {
					return fmt.Errorf("Issues.EditComment returned error: %v", err)
				}
				return nil
			}
		}
	}

	// Creates new comment
	_, _, err = client.Issues.CreateComment(ctx, owner, repo, prNumber, &github.IssueComment{Body: github.String(content)})
	if err != nil {
		return fmt.Errorf("Issues.CreateComment returned error: %v", err)
	}

	return nil
}
