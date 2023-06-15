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

func CommentReply() error {
	issueNumber := int(types.ActionConfigJSON.IssueOrPRNumber)
	fileName, _ := GetEnvFromAction("filename")
	comment, _ := GetEnvFromAction("comment")
	isReply, _ := GetEnvFromAction("isReply")
	if fileName == "" && comment == "" {
		return fmt.Errorf("filename or comment is empty")
	}

	if fileName != "" {
		fileContent, err := os.ReadFile(fileName)
		if err != nil {
			return err
		}
		comment = string(fileContent)
	}

	if isReply == "true" {
		replyBody, _, _ := unstructured.NestedString(types.ActionConfigJSON.Data, "comment", "body")
		replyBody = utils.QuoteReply(replyBody)
		comment = strings.Join([]string{replyBody, comment}, "\r\n\r\n")
	}

	owner, repo, err := getRepo()
	if err != nil {
		return err
	}
	logger.Info("repo:%s, issueNumber: %d", repo, issueNumber)
	logger.Debug("comment: %s", comment)
	ctx := context.Background()
	client := github_go.GithubClient(ctx)
	githubComment := &github.IssueComment{Body: github.String(comment)}
	_, _, err = client.Issues.CreateComment(ctx, owner, repo, issueNumber, githubComment)
	return err
}
