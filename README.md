## gh-rebot 使用说明文档

gh-rebot 是一个用 Go 语言编写的命令行工具，用于自动化 GitHub 操作。本文档将介绍如何安装和使用 gh-rebot。

### 安装

从 GitHub Actions 工作流中下载已编译的二进制文件。文件将保存在 dist/gh-rebot_linux_amd64_v1/gh-rebot 目录下。将其移动到适当的位置并设置可执行权限。

### 使用

gh-rebot 支持以下子命令：

1. version

```shell
gh-rebot version
```

2. changelog


```shell
gh-rebot changelog --reviews cuisongliu
```

3. comment

```shell
gh-rebot comment "${{github.event.comment.body}}"  ${{ github.event.issue.number }} ${{ github.repository_owner }}/gh-rebot
```

- GH_TOKEN：这是一个 GitHub Personal Access Token（个人访问令牌），用于对 GitHub API 进行身份验证。通常，它存储在仓库的 Secrets 中以保护敏感信息。在此例中，"${{ secrets.GH_PAT }}" 表示从仓库的 Secrets 中获取名为 GH_PAT 的变量值。
- SEALOS_SYS_TRIGGERS：这是一个字符串，其中包含以逗号分隔的 GitHub 用户名列表。这些用户名代表了允许触发 GitHub Actions 工作流的用户。在此例中，"sfggg" 是一个示例用户名，您需要用实际允许触发工作流的用户列表替换它。


1. **version**：表示配置文件的版本。在本例中，版本为 `v1`。
2. **debug**：布尔值（`true` 或 `false`），表示是否开启调试模式。在调试模式下，工具可能会输出更多的日志信息以帮助诊断问题。本例中，调试模式已开启。
3. **bot**：这部分包含与 GitHub 机器人相关的配置。
    - **allowOps**：字符串，包含允许操作此 GitHub 机器人的 GitHub 用户名。在本例中，只有 `cuisongliu` 用户允许操作此机器人。
    - **email**：字符串，表示机器人的电子邮件地址。本例中，电子邮件地址为 `sealos-ci-robot@sealos.io`。
    - **username**：字符串，表示机器人的 GitHub 用户名。本例中，用户名为 `sealos-ci-robot`。
4. **repo**：这部分包含与仓库相关的配置。
    - **org**：布尔值（`true` 或 `false`），表示仓库是否属于一个组织。本例中，仓库属于组织，值为 `true`。
    - **name**：字符串，表示主仓库的名称。在本例中，主仓库为 `labring-actions/sealos`。
    - **fork**：字符串，表示分支仓库的名称。在本例中，分支仓库为 `cuisongliu/sealos`。

### 示例

首先，创建名为 `.github/gh-bot.yml` 或 `.github/gh-bot.yaml` 的配置文件，并将以下内容粘贴到其中：

```yaml
version: v1
debug: true
bot:
  allowOps:
  - cuisongliu
  email: sealos-ci-robot@sealos.io
  username: sealos-ci-robot
repo:
  org: true
  name: labring-actions/sealos
  fork: cuisongliu/sealos

```

接下来，创建一个 GitHub Action 工作流文件。在您的项目仓库的 `.github/workflows` 目录中，创建一个名为 `gh_bot_action.yml` 的文件，并将以下内容粘贴到其中：

```yaml
name: Release Changelog
on:
   push:
      branches-ignore:
         - '**'
      tags:
         - '*'
env:
  # Common versions
  GH_TOKEN: "${{ secrets.GH_PAT }}"
  SEALOS_SYS_TRIGGER: "${{ github.event.sender.login }}"
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@master
      - name: Download gh-rebot
        run: |
          ##wget gh-rebot

      - name: Test sub cmd
        run: |
          gh-rebot changelog


```

