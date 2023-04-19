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
	"github.com/labring-actions/gh-rebot/pkg/types"
	"github.com/labring-actions/gh-rebot/pkg/utils"
	"k8s.io/apimachinery/pkg/util/json"
	"strings"
	"time"
)

type ActionOut struct {
	Conclusion string `json:"conclusion"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	URL        string `json:"url"`
	IsSuccess  bool   `json:"-"`
}

func CheckRelease(tagName string) (*ActionOut, error) {
	workflowOutput, _ := utils.RunCommandWithOutput(fmt.Sprintf(gitWorkflowCheck, types.GlobalsBotConfig.GetForkName(), types.GlobalsBotConfig.Release.Action, tagName), true)
	if workflowOutput == "" || strings.Contains(workflowOutput, "could not find any workflows named") {
		time.Sleep(5 * time.Second)
		return CheckRelease(tagName)
	}
	var out ActionOut
	if err := json.Unmarshal([]byte(workflowOutput), &out); err != nil {
		return nil, err
	}
	if out.Status == "completed" {
		if out.Conclusion == "success" {
			out.IsSuccess = true
			return &out, nil
		}
		out.IsSuccess = false
		return &out, nil
	} else {
		tt, err := time.ParseDuration(types.GlobalsBotConfig.Release.Retry)
		if err != nil {
			tt = time.Second * 20
		}
		logger.Debug("workflow release is in progress, please wait %s retry", tt.String())
		time.Sleep(tt)
		return CheckRelease(tagName)
	}
}
