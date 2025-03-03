package controller

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/ChamHerry/oshelper/consts"
	"github.com/ChamHerry/oshelper/utils"
	"github.com/gogf/gf/v2/frame/g"
)

// getCentosSystemDetailVersion 获取系统版本
func (s *Controller) getCentosSystemDetailVersion(config consts.SystemInfo) (ret string, err error) {
	command := "cat /etc/system-release"
	out, err := s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return config.VersionID, err
	}
	out = strings.TrimSpace(out)
	// 修改正则表达式以更准确地匹配版本号
	re := regexp.MustCompile(`(?i)([A-Za-z\s]+)\s+release\s+([\d.]+)\s*(?:\((.*?)\))?`)
	match := re.FindStringSubmatch(out)
	if len(match) > 0 {
		return match[2], nil
	}
	return config.VersionID, nil
}

// 获取已经安装的软件包列表
func (s *Controller) getCentosInstalledPackageList(ctx context.Context, in consts.GetInstalledPackageListParam) (out consts.GetInstalledPackageListResult, err error) {
	command := "rpm -qa"
	commandOut, err := s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return out, err
	}
	out.PackageList = strings.Split(strings.TrimSpace(commandOut), "\n")
	return out, nil
}

// getCentosPackageInfo 获取某个软件包的详细信息
func (s *Controller) getCentosPackageInfo(ctx context.Context, in consts.GetPackageInfoParam) (out consts.GetPackageInfoResult, err error) {
	var runCommandConfig consts.RunCommandConfig
	if strings.Contains(in.PackageName, "/") || strings.HasSuffix(in.PackageName, ".rpm") {
		runCommandConfig.Command = "rpm -qpi " + in.PackageName
	} else {
		runCommandConfig.Command = "rpm -qi " + in.PackageName
	}
	commandOut, err := s.RunCommand(runCommandConfig)
	if err != nil {
		return out, err
	}
	packageInfo := consts.PackageInfo{}
	// packageInfo.FullName = in.PackageName
	if strings.Contains(in.PackageName, "/") || strings.HasSuffix(in.PackageName, ".rpm") {
		splitList := strings.Split(in.PackageName, "/")
		packageInfo.FullName = splitList[len(splitList)-1]
	} else {
		packageInfo.FullName = in.PackageName + ".rpm"
	}
	// 解析命令输出
	lines := strings.Split(commandOut, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Name") && packageInfo.Name == "" {
			// 解析包名
			packageInfo.Name = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		} else if strings.HasPrefix(line, "Version") && packageInfo.Version == "" {
			// 解析版本号
			packageInfo.Version = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		} else if strings.HasPrefix(line, "Release") && packageInfo.ReleaseNumber == "" {
			// 解析发布号
			release := strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
			release = strings.Trim(release, "."+s.systemInfo.ID)
			releaseList := strings.Split(release, ".")
			if len(releaseList) == 1 {
				packageInfo.ReleaseNumber = release
				packageInfo.OS = ""
			} else {
				packageInfo.ReleaseNumber = strings.Join(releaseList[:len(releaseList)-1], ".")
				packageInfo.OS = releaseList[len(releaseList)-1]
			}
		} else if strings.HasPrefix(line, "Architecture") && packageInfo.Architecture == "" {
			// 解析架构
			packageInfo.Architecture = strings.TrimSpace(strings.SplitN(line, ":", 2)[1])
		}
	}
	// 检查是否所有字段都被正确解析
	if packageInfo.Name == "" || packageInfo.Version == "" || packageInfo.ReleaseNumber == "" {
		g.Log().Error(ctx, "packageInfo", packageInfo)
		g.Log().Error(ctx, "Failed to parse package details from rpm -qi output for:", in.PackageName)
		return out, fmt.Errorf("failed to parse package details for %s", in.PackageName)
	}
	out.PackageInfo = packageInfo
	return out, nil
}

// 获取某个软件包的文件列表
func (s *Controller) getCentosPackageFileList(ctx context.Context, in consts.GetPackageFileListParam) (out consts.GetPackageFileListResult, err error) {
	var runCommandConfig consts.RunCommandConfig
	if strings.Contains(in.PackageName, "/") || strings.HasSuffix(in.PackageName, ".rpm") {
		runCommandConfig.Command = "rpm -qpl " + in.PackageName
		out.PackageInfo.FullName = strings.Split(in.PackageName, "/")[len(strings.Split(in.PackageName, "/"))-1]
	} else {
		runCommandConfig.Command = "rpm -ql " + in.PackageName
		out.PackageInfo.FullName = in.PackageName + ".rpm"
	}
	commandOut, err := s.RunCommand(runCommandConfig)
	if err != nil {
		return out, err
	}
	out.PackageInfo.Files = strings.Split(strings.TrimSpace(commandOut), "\n")
	return out, nil
}

// installCentosPackage 安装软件包
func (s *Controller) installCentosPackages(ctx context.Context, in consts.InstallPackagesParam) (installPackageResult consts.InstallPackagesResult, err error) {
	installPackageResult = consts.InstallPackagesResult{}
	installPackageResult.IgnoredInstallPackageList = make([]consts.PackageInfo, 0)
	installPackageResult.SuccessfullyInstallPackageList = make([]consts.PackageInfo, 0)
	installPackageResult.FailedInstallPackageList = make([]consts.PackageInfo, 0)
	installPackageResult.Total = len(in.PackageList)
	// 获取已经安装的软件包列表
	installedPackageList, err := s.getCentosInstalledPackageList(ctx, consts.GetInstalledPackageListParam{})
	if err != nil {
		return installPackageResult, err
	}
	installedPackageListMap := make(map[string]consts.PackageInfo)
	for _, v := range installedPackageList.PackageList {
		packageInfo, err := utils.ParseRPMName(v)
		if err != nil {
			return installPackageResult, err
		}
		installedPackageListMap[v] = packageInfo
		g.Log().Debug(ctx, "packageInfo", packageInfo)
	}

	// for _, v := range in.PackageList {
	// 	g.Log().Debug(ctx, "v", v)
	// }

	return installPackageResult, nil
}
