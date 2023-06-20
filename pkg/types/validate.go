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

package types

import "fmt"

func (r *Config) validate() error {
	if r.Bot.Username == "" {
		return fmt.Errorf("bot username is required")
	}
	if r.Bot.Email == "" {
		return fmt.Errorf("bot email is required")
	}

	if r.Type == TypeAction {
		if r.GetRepoName() == "" {
			return fmt.Errorf("repo name is required")
		}
		if r.GetForkName() == "" {
			return fmt.Errorf("repo fork is required")
		}
		if r.Action.Release != nil {
			if r.Action.Release.ActionName == "" {
				return fmt.Errorf("release action is required")
			}
			if r.Action.Release.Retry == "" {
				return fmt.Errorf("release retry is required")
			}
		}
	}
	return nil
}

func (t *GithubVar) validate() error {
	if t.RunnerID == "" {
		return fmt.Errorf("error: GITHUB_RUN_ID is not set. Please set the GITHUB_RUN_ID environment variable")
	}
	if t.RepoFullName == "" {
		return fmt.Errorf("error: not found repository.full_name in github event")
	}
	return nil
}

func ValidateIssueOrPRNumber() error {
	if ActionConfigJSON.IssueOrPRNumber == 0 {
		return fmt.Errorf("error: not found issue.number or pull_request.number in github event")
	}
	return nil
}
