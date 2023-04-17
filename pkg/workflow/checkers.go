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
	"github.com/labring-actions/gh-rebot/pkg/config"
	"github.com/labring-actions/gh-rebot/pkg/gh"
	"github.com/labring-actions/gh-rebot/pkg/utils"
	"github.com/pkg/errors"
)

func checkPermission(extras []string) error {
	tigger := gh.GlobalsGithubVar.SenderOrCommentUser
	if tigger == "" {
		return errors.New("Error: github sender or workflow is empty.")
	}
	ops := config.GlobalsConfig.GetBotAllowOps()
	ops = append(ops, extras...)

	if len(ops) == 0 {
		return errors.New("Error: no has permission users to trigger this action.")
	}

	if utils.In(ops, "all") {
		return nil
	}

	if !utils.In(ops, tigger) {
		return errors.New("Error: no has permission to trigger this action.")
	}
	return nil
}
