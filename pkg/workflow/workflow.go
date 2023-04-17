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
	"github.com/cuisongliu/logger"
	"github.com/labring-actions/gh-rebot/pkg/config"
	"github.com/labring-actions/gh-rebot/pkg/gh"
	"github.com/labring-actions/gh-rebot/pkg/utils"
	"strings"
)

type workflow struct {
	Body string
}

func (c *workflow) Release() error {
	if checkPermission(config.GlobalsConfig.Release.AllowOps) != nil {
		return sendMsgToIssue("permission_error")
	}
	data := strings.Split(c.Body, " ")
	if len(data) == 2 && utils.ValidateVersion(data[1]) {

		return sendMsgToIssue("success")
	} else {
		logger.Error("command format is error: %s ex. /{prefix}_release {tag}", c.Body)
		return sendMsgToIssue("format_error")
	}
}

func (c *workflow) Changelog() error {
	if checkPermission(config.GlobalsConfig.Changelog.AllowOps) != nil {
		return sendMsgToIssue("permission_error")
	}
	data := strings.Split(c.Body, " ")
	if len(data) == 1 {

		return sendMsgToIssue("success")
	} else {
		logger.Error("command format is error: %s ex. /{prefix}_changelog", c.Body)
		return sendMsgToIssue("format_error")
	}
}

func sendMsgToIssue(msgKey string) error {
	msg := config.GlobalsConfig.GetMessage(msgKey)
	return gh.SendMsgToIssue(msg)
}

func NewWorkflow(body string) Interface {
	return &workflow{Body: body}
}
