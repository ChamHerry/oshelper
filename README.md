# OSHelper

## 项目简介

OSHelper是一个开源项目，旨在通过SSH连接远程服务器，获取系统信息、管理软件包以及执行文件操作。该项目使用Go语言开发，主要功能包括获取系统详细版本信息、安装和管理软件包、下载和删除文件等。

## 功能特性

- **获取系统信息**: 支持获取Linux和Mac系统的详细信息，包括系统版本、架构、内核版本等。
- **软件包管理**: 支持CentOS和Debian系统的软件包安装、信息获取和文件列表获取。
- **文件操作**: 支持远程文件的下载和删除操作。
- **网络测试**: 提供简单的网络连通性测试功能。

## 安装与使用

1. 安装OSHelper：
   ```bash
   go get github.com/ChamHerry/oshelper
   ```
2. 在你的Go项目中导入并使用OSHelper：
   ```go
   import "github.com/ChamHerry/oshelper"
   ```

3. 运行示例：
   ```bash
   package main
   
   import (
   	"context"
   	"github.com/ChamHerry/oshelper/consts"
   	"github.com/ChamHerry/oshelper/controller"
   	"github.com/gogf/gf/v2/frame/g"
   )
   
   func main() {
   	ctx := context.Background()
   	SystemController, _ := controller.NewController(consts.SSHConfig{})
   	runCommand, err := SystemController.RunCommand(consts.RunCommandConfig{Command: "ls"})
   	if err != nil {
   		return
   	}
   	g.Log().Info(ctx, runCommand)
   }
   ```

## 贡献

欢迎贡献！请阅读[贡献指南](CONTRIBUTING.md)以了解如何参与项目。

## 许可证

该项目基于GPL 3.0许可证进行许可。详细信息请参阅[LICENSE](LICENSE)文件。
