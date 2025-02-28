package controller

import (
	"context"
	"oshelper/consts"
	"path/filepath"
	"regexp"
	"strings"
)

func (s *Controller) getDebianSystemDetailVersion(config consts.SystemInfo) (ret string, err error) {
	command := "cat /etc/lsb-release"
	out, err := s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return config.VersionID, err
	}
	out = strings.TrimSpace(out)
	re := regexp.MustCompile(`(?i)DISTRIB_RELEASE=([\d.]+)`)
	match := re.FindStringSubmatch(out)
	if len(match) > 0 {
		return match[1], nil
	}
	return config.VersionID, nil
}

// 获取已经安装的软件包列表
func (s *Controller) getDebianInstalledPackageList(ctx context.Context, in consts.GetInstalledPackageListParam) (out consts.GetInstalledPackageListResult, err error) {
	command := "dpkg --get-selections|awk -F ' ' '{print $1}'|awk -F ' ' '{print $1}'"
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

// getDebianPackageInfo 获取某个软件包的详细信息
func (s *Controller) getDebianPackageInfo(ctx context.Context, in consts.GetPackageInfoParam) (out consts.GetPackageInfoResult, err error) {
	var commandOut string
	fullName := ""
	if strings.Contains(in.PackageName, "http") {
		// 下载文件
		downloadFileResult, err := s.DownloadFile(ctx, consts.DownloadFileParam{
			URL: in.PackageName,
		})
		if err != nil {
			return out, err
		}
		// 获取文件信息
		commandOut, err = s.getDebianPackageInfoByFile(ctx, consts.GetPackageInfoParam{
			PackageName: downloadFileResult.FilePath,
		})
		if err != nil {
			return out, err
		}
		// 删除文件
		defer s.DeleteFile(ctx, consts.DeleteFileParam{
			FilePath: filepath.Dir(downloadFileResult.FilePath),
		})
		// 获取文件名
		splitList := strings.Split(in.PackageName, "/")
		fullName = splitList[len(splitList)-1]
	} else if strings.HasSuffix(in.PackageName, ".deb") {
		commandOut, err = s.getDebianPackageInfoByFile(ctx, in)
		// 获取文件名
		splitList := strings.Split(in.PackageName, "/")
		fullName = splitList[len(splitList)-1]
	} else {
		var runCommandConfig consts.RunCommandConfig
		runCommandConfig.Command = "dpkg -s " + in.PackageName
		commandOut, err = s.RunCommand(runCommandConfig)
	}
	if err != nil {
		return out, err
	}
	packageInfo := consts.PackageInfo{}
	// 解析命令输出
	lines := strings.Split(commandOut, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		splitList := strings.Split(line, ":")
		if len(splitList) < 2 {
			continue
		}
		key := strings.TrimSpace(splitList[0])
		value := strings.TrimSpace(splitList[1])
		switch key {
		case "Package":
			packageInfo.Name = value
		case "Version":
			packageInfo.Version = value
		case "Architecture":
			packageInfo.Architecture = value
		}
	}
	// fullName 是包名_版本_架构
	if fullName != "" {
		packageInfo.FullName = fullName
	} else {
		packageInfo.FullName = packageInfo.Name + "_" + packageInfo.Version + "_" + packageInfo.Architecture + ".deb"
	}
	out.PackageInfo = packageInfo
	return out, nil
}

// getDebianPackageInfoByFile 通过文件获取软件包的详细信息
func (s *Controller) getDebianPackageInfoByFile(ctx context.Context, in consts.GetPackageInfoParam) (out string, err error) {
	command := "dpkg --info " + in.PackageName
	return s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})

}

// 获取某个软件包的文件列表
func (s *Controller) getDebianPackageFileList(ctx context.Context, in consts.GetPackageFileListParam) (out consts.GetPackageFileListResult, err error) {
	var commandOut string
	if strings.Contains(in.PackageName, "http") {
		// 下载文件
		downloadFileResult, err := s.DownloadFile(ctx, consts.DownloadFileParam{
			URL: in.PackageName,
		})
		if err != nil {
			return out, err
		}
		// 删除文件
		defer s.DeleteFile(ctx, consts.DeleteFileParam{
			FilePath: filepath.Dir(downloadFileResult.FilePath),
		})
		// 获取文件信息
		return s.getDebianPackageFileListByFile(ctx, consts.GetPackageFileListParam{
			PackageName: downloadFileResult.FilePath,
		})

	} else if strings.Contains(in.PackageName, ".deb") {

		return s.getDebianPackageFileListByFile(ctx, in)
	}
	command := "dpkg -L " + in.PackageName
	commandOut, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return out, err
	}
	var outList []string
	lines := strings.Split(strings.TrimSpace(commandOut), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == "/." {
			continue
		}
		outList = append(outList, line)
	}
	// 获取文件名
	packageInfo, err := s.getDebianPackageInfo(ctx, consts.GetPackageInfoParam{
		PackageName: in.PackageName,
	})
	if err != nil {
		return out, err
	}
	out.PackageInfo.Files = outList
	out.PackageInfo.FullName = packageInfo.FullName
	return out, nil
}

// getDebianPackageFileListByFile 通过文件获取软件包的文件列表
func (s *Controller) getDebianPackageFileListByFile(ctx context.Context, in consts.GetPackageFileListParam) (out consts.GetPackageFileListResult, err error) {
	command := "dpkg --contents " + in.PackageName
	commandOut, err := s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return out, err
	}
	outList := make([]string, 0)
	lines := strings.Split(commandOut, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		splitList := strings.Split(line, " ")
		if len(splitList) < 2 {
			continue
		}
		var filePath string
		if splitList[0][0] == 'l' {
			filePath = splitList[len(splitList)-3]
		} else {
			filePath = splitList[len(splitList)-1]
		}
		filePath = strings.TrimSpace(filePath)

		if strings.HasPrefix(filePath, "./") {
			filePath = filePath[1:]
		}
		if filePath == "/" {
			continue
		}
		filePath = strings.TrimRight(filePath, "/")
		// g.Log().Debug(ctx, "splitList:", splitList)
		outList = append(outList, filePath)
	}
	out.PackageInfo.Files = outList
	// 获取文件名
	splitList := strings.Split(in.PackageName, "/")
	out.PackageInfo.FullName = splitList[len(splitList)-1]
	return out, nil
}
