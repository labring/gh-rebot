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

const (
	authStatus        = "gh auth status"
	disablePrompt     = "gh config set prompt disabled"
	configEmail       = "cd /tmp/fork-sealos-repo && git config user.email %s"
	configUser        = "cd /tmp/fork-sealos-repo && git config user.name %s"
	forkRepo          = "gh repo view %s >/dev/null 2>&1 || (echo '仓库不存在，正在fork仓库...' && gh repo fork %s %s)"
	syncRepo          = "cd /tmp/fork-sealos-repo && git fetch fork && git rebase fork/main && git push -f origin main && git fetch --tags fork && git push -f --tags origin"
	cloneRepo         = "rm -rf /tmp/fork-sealos-repo && && gh repo clone %s /tmp/fork-sealos-repo"
	setToken          = "cd /tmp/fork-sealos-repo && git remote set-url origin https://%s:%s@github.com/%s.git"
	newBranch         = "cd /tmp/fork-sealos-repo && git checkout -b %s"
	generateChangelog = "cd /tmp/fork-sealos-repo && bash %s"
	gitStatus         = "cd /tmp/fork-sealos-repo && git status"
	gitAdd            = "cd /tmp/fork-sealos-repo && git add -A"
	gitCommit         = "cd /tmp/fork-sealos-repo && git commit -am '%s' -s"
	gitRelease        = "cd /tmp/fork-sealos-repo && git describe --abbrev=0 --tags"
	gitPush           = "cd /tmp/fork-sealos-repo && git push -u origin %s"
	gitPR             = "cd /tmp/fork-sealos-repo && gh pr create --title 'docs: Automated Changelog Update for %s' --body '🤖 add release changelog using rebot' --reviewer %s "
	gitTag            = "cd /tmp/fork-sealos-repo && git fetch --tags && git tag -l"
	gitNewTag         = "cd /tmp/fork-sealos-repo && git tag %s"
	gitAddRemote      = "cd /tmp/fork-sealos-repo && git remote add fork https://%s:%s@github.com/%s.git"
	gitPushRemote     = "cd /tmp/fork-sealos-repo && git push -f fork %s"
	gitWorkflowCheck  = "cd /tmp/fork-sealos-repo && gh run list --repo %s -w %s -b %s --json name,status,url,conclusion -q '.[0]'"
)
