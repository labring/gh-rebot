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

package workflow

import (
	"bytes"
	"github.com/labring/gh-rebot/pkg/template"
	"github.com/labring/gh-rebot/pkg/types"
	"github.com/labring/gh-rebot/pkg/utils"
	"github.com/pkg/errors"
)

type sender struct {
	Body  string
	Error string
}

func (s *sender) sendMsgToIssue(msgKey string, actionURL ...string) error {
	msg := types.GlobalsBotConfig.GetMessage(msgKey)
	v, b, _ := template.TryParse(msg)
	if b {
		out := bytes.NewBuffer(nil)
		_ = v.Execute(out, map[string]interface{}{
			"Body":  s.Body,
			"Error": s.Error,
		})
		msg = out.String()
	}
	return SendMsgToIssue(msg, actionURL...)
}

func (s *sender) sendCommentMsgToIssue(msg string) error {
	return SendCustomizeMsgToIssue(msg)
}

func SendMsgToIssue(msg string, actionURL ...string) error {
	tpl, ok, _ := template.TryParse(`gh issue comment {{.IssueOrPRNumber}} --body "{{.Msg}} <br/>See: <br/>- {{.GetRunnerURL}}{{range .ActionURLs}}<br/>- {{.}}{{end}}" --repo {{.SafeRepo}}`)
	if ok {
		out := bytes.NewBuffer(nil)
		_ = tpl.Execute(out, map[string]interface{}{
			"IssueOrPRNumber": types.ActionConfigJSON.IssueOrPRNumber,
			"Msg":             msg,
			"GetRunnerURL":    types.ActionConfigJSON.GetRunnerURL(),
			"SafeRepo":        types.ActionConfigJSON.SafeRepo,
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
			"IssueOrPRNumber": types.ActionConfigJSON.IssueOrPRNumber,
			"Msg":             msg,
			"SafeRepo":        types.ActionConfigJSON.SafeRepo,
		})
		return utils.RunCommand("bash", "-c", out.String())
	}

	return errors.New("template parse error")
}
