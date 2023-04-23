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
	"github.com/labring-actions/gh-rebot/pkg/gh"
	"github.com/labring-actions/gh-rebot/pkg/types"
	"strings"
)

type changelog struct {
	*workflow
}

func (c *changelog) Run() error {
	if checkPermission(types.GlobalsBotConfig.Changelog.AllowOps) != nil {
		return c.sender.sendMsgToIssue("permission_error")
	}
	data := strings.Split(c.Body, " ")
	if len(data) == 1 {
		err := gh.Changelog(types.GlobalsBotConfig.Changelog.Reviewers)
		if err != nil {
			c.sender.Error = err.Error()
			return c.sender.sendMsgToIssue("changelog_error")
		}
		return c.sender.sendMsgToIssue("success")
	} else {
		logger.Error("command format is error: %s ex. /{prefix}_changelog", c.Body)
		return c.sender.sendMsgToIssue("format_error")
	}
}

func (c *changelog) Comment() string {
	if types.GlobalsBotConfig.GetPrefix() == "/" {
		return "/changelog"
	}
	return strings.Join([]string{types.GlobalsBotConfig.GetPrefix(), "changelog"}, types.GlobalsBotConfig.GetSpe())
}

func NewChangelog(body string) Interface {
	return &changelog{workflow: newWorkflow(body)}
}
