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
	"github.com/labring/gh-rebot/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"os"
	"strings"
)

// IssueComment is a action to comment on PR
func IssueComment() error {
	fileName, _ := GetEnvFromAction("filename")
	commentBody, _ := GetEnvFromAction("comment")
	if fileName == "" && commentBody == "" {
		return fmt.Errorf("filename or comment is empty")
	}
	replaceTag, _ := GetEnvFromAction("replace_tag")
	prNumber := int(types.ActionConfigJSON.IssueOrPRNumber)

	if fileName != "" {
		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		commentBody = string(fileContent)
	}
	hiddenReplace := fmt.Sprintf("<!-- %s -->", replaceTag)
	if replaceTag != "" {
		commentBody = commentBody + "\n" + hiddenReplace
	}

	isReply, _ := GetEnvFromAction("isReply")
	if isReply == "true" {
		replyBody, _, _ := unstructured.NestedString(types.ActionConfigJSON.Data, "comment", "body")
		replyBody = utils.QuoteReply(replyBody)
		commentBody = strings.Join([]string{replyBody, commentBody}, "\r\n\r\n")
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
	// Checks existing comments, edits if match found
	if replaceTag != "" {
		for _, comment := range comments {
			if comment.Body != nil && comment.ID != nil {
				if hiddenReplace != "" && strings.LastIndex(*comment.Body, hiddenReplace) != -1 {
					_, _, err = client.Issues.EditComment(ctx, owner, repo, *comment.ID, &github.IssueComment{Body: github.String(commentBody)})
					if err != nil {
						return fmt.Errorf("Issues.EditComment returned error: %v", err)
					}
					return nil
				}
			}
		}
	}

	// Creates new comment
	_, _, err = client.Issues.CreateComment(ctx, owner, repo, prNumber, &github.IssueComment{Body: github.String(commentBody)})
	if err != nil {
		return fmt.Errorf("Issues.CreateComment returned error: %v", err)
	}

	return nil
}
