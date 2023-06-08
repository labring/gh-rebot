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
	"github.com/cuisongliu/logger"
	"github.com/labring/gh-rebot/pkg/types"
	"github.com/labring/gh-rebot/pkg/workflow"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"strings"
)

// CommentEngine is the main function for comment action
func CommentEngine() error {
	comment, _, _ := unstructured.NestedString(types.ActionConfigJSON.Data, "comment", "body")
	cmds := strings.Split(comment, "\n")
	for _, t := range cmds {
		logger.Debug("cmds: ", strings.TrimSpace(t))
		wfs := make([]workflow.Interface, 0)
		if types.GlobalsBotConfig.Release != nil {
			wfs = append(wfs, workflow.NewRelease(strings.TrimSpace(t)))
		}
		used := 0
		for _, wf := range wfs {
			if strings.HasPrefix(t, wf.Comment()) {
				if err := wf.Run(); err != nil {
					return err
				}
				used++
			}
		}
		if used == 0 {
			logger.Warn("not support command: ", t)
		}
	}
	return nil
}
