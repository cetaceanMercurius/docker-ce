# Docker CE

:warning: **这个存储库现已经废弃, 将被归档** (Docker CE 本身并**不是**废弃的) :warning:

从Docker 20.10版本开始，Docker Engine 和 Docker CLI 的软件包是直接从它们各自的源码库构建的，而不是从这个库。

实际上, 这意味着:
1. 这个资源库不再是Docker CE构建的 "真理之源 "了。
2. Docker CLI构建的提交SHA和标签将来自于 [docker/cli](https://github.com/docker/cli) 存储库, 
   而Docker引擎的提交SHA和标签将来自于 [moby/moby](https://github.com/moby/moby) 仓库。
4. Engine, CLI和打包的发布分支将被维护在他们各自的存储库。
4. 这个资源库将停止更新，将来会被归档。
5. Changelog 现在是 [Release Notes](https://docs.docker.com/engine/release-notes/).
6. 该版本库的 "master "分支将在版本库归档时被清空。

## 描述

这个仓库存放了Docker CE产品的开源组件。`master`分支的作用是定期统一上游的组件。Long-lived 的发布分支承载着产品版本中在生命周期内的代码。

该资源库由 Docker 独家维护。

## Issues

有单独的问题跟踪库，用于终端用户的Docker CE产品，专门用于某个平台。
找到你的问题或为你使用的平台提交一个新问题:

* https://github.com/docker/for-linux
* https://github.com/docker/for-mac
* https://github.com/docker/for-win

## 提交pull请求

这个版本库不直接接受 [components](组件) 目录下的文件的 PR。
要对组件目录下的文件进行贡献，请参见 [CONTRIBUTING.md](CONTRIBUTING.md#submitting-pull-requests) 。

## 统一上游资源

`master`分支是使用 [moby-components]https://github.com/shykes/moby-extras/blob/master/cmd/moby-components) 工具
将不同上游git仓库的组件改编成统一的目录结构。

你可以在 [components.conf](components.conf) 文件中查看上游的git repos。
每个组件都被隔离到[components](components)目录下的自己的目录中。

该工具将在适当的路径内导入每个组件的git历史。

例如，这显示了一个提交被导入到组件`engine`, 
从 [moby/moby@a27b4b8](https://github.com/moby/moby/commit/a27b4b8cb8e838d03a99b6d2b30f76bdaf2f9e5d) 到 `components/engine`目录。

```
commit 5c70746915d4589a692cbe50a43cf619ed0b7152
Author: Andrea Luzzardi <aluzzardi@gmail.com>
Date:   Sat Jan 19 00:13:39 2013

    Initial commit
    Upstream-commit: a27b4b8cb8e838d03a99b6d2b30f76bdaf2f9e5d
    Component: engine

 components/engine/container.go       | 203 ++++++++++++++++++++++++++++...
 components/engine/container_test.go  | 186 ++++++++++++++++++++++++++++...
 components/engine/docker.go          | 112 ++++++++++++++++++++++++++++...
 components/engine/docker_test.go     | 175 ++++++++++++++++++++++++++++...
 components/engine/filesystem.go      |  52 ++++++++++++++++++++++++++++...
 components/engine/filesystem_test.go |  35 +++++++++++++++++++++++++++
 components/engine/lxc_template.go    |  94 ++++++++++++++++++++++++++++...
 components/engine/state.go           |  48 ++++++++++++++++++++++++++++...
 components/engine/utils.go           | 115 ++++++++++++++++++++++++++++...
 components/engine/utils_test.go      | 126 ++++++++++++++++++++++++++++...
 10 files changed, 1146 insertions(+)
```

## 对 "master "分支的更新

新功能的主要开发应针对上游的git仓库。
这个仓库的 "master "分支将定期从上游拉入新的变化，以提供一个整合点。

## 释放的分支

当Docker CE开始发布时，将从`master`创建一个新的分支。
分支名称将是`YY.MM`，代表基于时间的产品发布版本，例如`17.06`。

## 向发布分支添加修复程序

注意：每次提交的修正都应该只影响一个组件目录下的文件。

### 上游可用的修复方法

A PR cherry-picking the necessary commits should be created against
the release branch. If the the cherry-pick cannot be applied cleanly,
the logic of the fix should be ported manually.

### 暂未修复

首先为发布分支创建带有修正的 PR。
一旦修复被合并，一定要把修复移植到相应的上游git repo。

## 发布标签

每个候选版本（RC）和通用版本（GA）都会有一个git标签。
该标签将只指向发布分支上的提交。
