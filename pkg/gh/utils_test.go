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
	"github.com/labring-actions/gh-rebot/pkg/types"
	"testing"
)

func Test_checkRemoteTagExists(t *testing.T) {
	types.GlobalsGithubVar = new(types.GithubVar)
	types.GlobalsGithubVar.IssueOrPRNumber = 1
	types.GlobalsGithubVar.SafeRepo = "cuisongliu/sealos"
	types.GlobalsGithubVar.RunnerID = "12345445"
	SendMsgToIssue("default", "https://baidu.com", "https://sealos.io")
}
