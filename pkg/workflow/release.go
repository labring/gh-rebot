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
	"fmt"
	"github.com/cuisongliu/logger"
	"github.com/labring/gh-rebot/pkg/gh"
	"github.com/labring/gh-rebot/pkg/types"
	"github.com/labring/gh-rebot/pkg/utils"
	"strings"
)

type release struct {
	*workflow
}

func (c *release) Run() error {
	if checkPermission(types.GlobalsBotConfig.Action.Release.AllowOps) != nil {
		return c.sender.sendMsgToIssue("permission_error")
	}
	data := strings.Split(c.Body, " ")
	if len(data) == 2 && utils.ValidateVersion(data[1]) {
		err := gh.Tag(data[1])
		if err != nil {
			c.sender.Error = err.Error()
			return c.sender.sendMsgToIssue("release_error")
		}
		action, err := gh.CheckRelease(data[1])
		if err != nil {
			return err
		}
		if !action.IsSuccess {
			c.sender.Error = fmt.Sprintf("release action status is %s,action conclusion is %s", action.Status, action.Conclusion)
			return c.sender.sendMsgToIssue("release_error", action.URL)
		}
		if err = c.sender.sendMsgToIssue("success", action.URL); err != nil {
			return err
		}
		return nil
	} else {
		logger.Error("command format is error: %s ex. /{prefix}_release {tag}", c.Body)
		return c.sender.sendMsgToIssue("format_error")
	}
}

func (c *release) Comment() string {
	if types.GlobalsBotConfig.GetPrefix() == "/" {
		return "/release"
	}
	return strings.Join([]string{types.GlobalsBotConfig.GetPrefix(), "release"}, types.GlobalsBotConfig.GetSpe())
}

func NewRelease(body string) Interface {
	return &release{workflow: newWorkflow(body)}
}
