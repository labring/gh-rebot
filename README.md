# 项目已经废弃，请移步 [robot](https://github.com/labring/robot)

# gh-rebot 项目说明文档

gh-rebot 是一个针对 sealos 项目的 GitHub rebot，用于自动执行一些常见操作，如发布新版本等。本文档将介绍该 rebot 的配置文件，并提供相应的使用指南。

## 配置文件

下面是 gh-rebot 项目的配置文件：

```yaml
version: v1
debug: true
bot:
  prefix: /sealos
  action:
    printConfig: true
    release:
      retry: 15s
      action: Release
      allowOps:
        - cuisongliu
  spe: _
  allowOps:
    - sealos-ci-robot
    - sealos-release-robot
  email: sealos-ci-robot@sealos.io
  username: sealos-ci-robot
repo:
  org: true
  name: labring/sealos
  fork: cuisongliu/sealos

message:
  success: |
    🤖 says: The action {{.Body}} finished successfully 🎉
  format_error: |
    🤖 says: ‼️ The action format error, please check the format of this action.
  permission_error: |
    🤖 says: ‼️ The action no has permission to trigger.
  release_error: |
    🤖 says: ‼️ The action release error.

```

### 配置文件详解

- `version` - 版本标识，当前为 v1。
- `debug` - 是否开启调试模式，设置为 true 时开启。
- `action` \- action配置。
   - `printConfig` - 是否打印配置信息，设置为 true 时打印。
   - `release` \- 发布配置。
       - `retry` - 重试间隔，例如：15s。
       - `action` - 执行动作，例如：Release。
       - `allowOps` - 允许触发发布操作的用户名列表。
- `bot` \- 机器人配置。
   - `prefix` - 机器人命令前缀，用于识别命令。默认值 `/`,如果设置为`/` 则 `spe` 失效。命令为`/release`
   - `spe` - 机器人命令分隔符，用于识别命令。默认值 `_`
   - `allowOps` - 允许操作的用户名列表。
   - `email` - 机器人邮箱。
   - `username` - 机器人用户名。
- `repo` \- 仓库配置。
   - `org` - 是否为组织仓库，设置为 true 时表示是组织仓库。
   - `name` - 仓库名称。
   - `fork` - fork 的仓库名称。
- `message` \- 消息配置。
   - `success` - 成功消息模板。
   - `format_error` - 格式错误消息模板。
   - `permission_error` - 权限错误消息模板。
   - `release_error` - 发布错误消息模板。

## 使用文档

使用 gh-rebot 时，需要遵循以下步骤：

1. 将配置文件添加到项目的`.github`目录` gh-bot.yml `文件。
2. 确保配置文件中的用户名、仓库名称等信息与实际情况相符。
3. 根据配置文件中的命令前缀（如本例中的 `/sealos`）在 GitHub 仓库的 issue 或 PR 中发表评论，以触发相应的操作。

### 变更日志操作

之前的操作已经废弃，使用 https://github.com/labring/sealos/blob/d528d6be713b9b9cf92169e5822d354d29fffb9d/.github/workflows/release.yml#L72



### 发布操作

如果需要发布新版本，请在 issue 或 PR 中使用以下命令：

```
/sealos_release
```

### 错误处理

根据配置文件中的消息模板，gh-rebot 会在执行操作过程中遇到错误时返回相应的提示消息。例如：

- 格式错误：‼️ 机器人说：操作格式错误，请检查此操作的格式。
- 权限错误：‼️ 机器人说：操作无权限触发。
- 发布错误：‼️ 机器人说：操作发布错误。

在遇到错误时，请根据提示信息进行相应的调整。


### 如何使用Action

```yaml
- name: Gh Rebot for Sealos
  uses: labring/gh-rebot@v0.0.6-rc6
  with:
    version: v0.0.6-rc6
  env:
    SEALOS_TYPE: "/comment"
    GH_TOKEN: "${{ secrets.GH_PAT }}"
```
**版本支持**: 

- [x] 支持release
  > 目标分支为`release-v1.2`，如果没有则默认为`main`分支,该功能v0.0.7-rc1支持
  
  `SEALOS_TYPE: "/comment"` # 评论触发
  example:
  ```markdown
   /release v1.2.1
   /release v1.2.3 release-v1.2 
  ```
  
- [x] 支持文本替换回复
  > 该功能v0.0.8-rc2 支持 (升级后新增了，SEALOS_COMMENT、SEALOS_ISREPLY)
  - `SEALOS_TYPE: "issue_comment"` # PR文本替换回复
  - `SEALOS_FILENAME: "README.md"` # PR文本替换回复文件位置
  - `SEALOS_COMMENT: "/xxxx"` # comment的内容
  - `SEALOS_REPLACE_TAG: "TAG"` # 寻找标记，根据这个标记进行替换
  - `SEALOS_ISREPLY: "true"` # 是否回复，根据当前的comment的内容追加

- [x] issue自动创建
  > 该功能v0.0.8-rc1支持

  入参:

  - `SEALOS_TYPE: "issue_renew"` # issue自动创建，支持回复comment
  - `SEALOS_ISSUE_TITLE: "dxxxx"` # issue的title
  - `SEALOS_ISSUE_BODY: "xxxx"` # issue内容
  - `SEALOS_ISSUE_BODYFILE: "README.md"` # issue内容如果多可以写文件
  - `SEALOS_ISSUE_LABEL: "dayly-report"` # 新增issue的label
  - `SEALOS_ISSUE_TYPE: "day"` # day和week , 会在titles上自动加上日期,day代表一天一个issue会关闭之前的issue，week以此类推
  - `SEALOS_ISSUE_REPO`: "sealos/sealos" # issue创建的仓库
  - `SEALOS_COMMENT_BODY`: "xxxx" # issue创建后的comment内容
  - `SEALOS_COMMENT_BODYFILE`: "xxxx" # issue创建后的comment内容如果多可以写文件
  
  返回参数：
  
  - env.SEALOS_ISSUE_NUMBER # issue的number

## Roadmap

- [ ] 支持label操作
- [ ] 支持里程碑操作
- [ ] 支持pr的code review操作
- [ ] 支持pr的merge操作
- [ ] 支持pr的close操作
- [ ] 支持pr的reopen操作
- [ ] 支持pr的comment操作
- [ ] 支持pr和issue的assign操作
