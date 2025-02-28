package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os/exec"
	"oshelper/consts"
	"oshelper/utils"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

type Controller struct {
	client     *ssh.Client
	systemInfo *consts.SystemInfo
}

// NewController 创建一个 Controller
func NewController(config consts.SSHConfig) (controller *Controller, err error) {
	controller = &Controller{}
	// 如果IP地址为空，则直接返回
	if config.IPAddress == "" {
		return controller, nil
	}
	if config.Port == 0 {
		config.Port = consts.DefaultSSHConfig.Port
	}
	if config.Username == "" {
		config.Username = consts.DefaultSSHConfig.Username
	}
	conn, err := net.DialTimeout("tcp", config.IPAddress+":"+strconv.Itoa(config.Port), consts.DefaultDialTimeout)
	if err != nil {
		return controller, err
	}
	defer conn.Close()

	sshConfig := &ssh.ClientConfig{
		User: config.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	ServerAddress := config.IPAddress + ":" + strconv.Itoa(config.Port)
	client, err := ssh.Dial("tcp", ServerAddress, sshConfig)
	if err != nil {
		return controller, err
	}
	controller.client = client
	controller.systemInfo, err = controller.GetSystemInfo(context.Background())
	if err != nil {
		return controller, err
	}

	// controller.systemInfo = &tempSystemInfo
	return controller, nil
}

// Close 关闭 SSH 连接
func (s *Controller) Close() {
	if s.client != nil {
		_ = s.client.Close()
	}
}

// RunLocalCommand 执行本地命令
func (s *Controller) RunLocalCommand(config consts.RunCommandConfig) (ret string, err error) {
	// 创建一个命令，运行 "ls -l" 命令
	// g.Log().Debug(context.Background(), "RunLocalCommand:", config.Command)
	cmd := exec.Command("bash", "-c", config.Command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		config.RunCommandFailedCounts += 1
		err = fmt.Errorf("failed to run command: %s", strings.TrimSpace(stderr.String()))
		if config.RunCommandFailedCounts >= consts.DefaultRunCommandFailedCounts {
			// g.Log().Debug(context.Background(), "RunLocalCommand err:", err)
			return ret, err
		}
		//return s.RunCommand(config)
		return s.RunLocalCommand(config)
	}
	// g.Log().Debugf(context.Background(), "RunLocalCommand stdout:%s", stdout.String())
	// g.Log().Debugf(context.Background(), "RunLocalCommand stderr:%s", stderr.String())
	ret, errStr := strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String())
	// g.Log().Debugf(context.Background(), "RunLocalCommand ret:%s, errStr:%s", ret, errStr)
	if errStr != "" {
		if strings.Contains(errStr, "警告") || strings.Contains(errStr, "warning") {
			return ret, nil
		}
		return ret, fmt.Errorf("failed to run command: %s", errStr)
	}
	// g.Log().Debugf(context.Background(), "Finish RunLocalCommand ret:%s", ret)
	return ret, nil
}

// RunCommand 执行远程命令
func (s *Controller) RunCommand(config consts.RunCommandConfig) (ret string, err error) {
	//ctx := s.ctx
	//g.Log().Debug(ctx, "Run Command:"+config.Command+", Counts:"+strconv.Itoa(config.RunCommandFailedCounts+1))
	if s.client == nil {
		return s.RunLocalCommand(config)
	}
	session, err := s.client.NewSession()
	if err != nil {
		return ret, fmt.Errorf("failed to create session: %w", err)
	}
	// 显示包函数忽略报错
	defer func(session *ssh.Session) {
		_ = session.Close()
	}(session)
	//var stdoutBuf bytes.Buffer
	var stdout, stderr bytes.Buffer
	// 将标准输出和错误重定向到缓冲区
	session.Stdout = &stdout
	session.Stderr = &stderr
	//session.Stdout = &stdoutBuf
	err = session.Run(config.Command)
	// 标准输出
	ret = stdout.String()
	if err != nil {
		// g.Log().Debug(context.Background(), "err:", err)
		// g.Log().Debug(context.Background(), "stderr.String():", stderr.String())
		// g.Log().Debug(context.Background(), "stdout.String():", ret)
		config.RunCommandFailedCounts += 1
		if stderr.String() != "" {
			err = fmt.Errorf("failed to run command: %s", strings.TrimSpace(stderr.String()))
		} else {
			err = fmt.Errorf("failed to run command: %s", strings.TrimSpace(ret))

		}
		if config.RunCommandFailedCounts >= consts.DefaultRunCommandFailedCounts {
			//g.Log().Debug(ctx, "Run Command:"+config.Command+", Counts:"+strconv.Itoa(config.RunCommandFailedCounts+1))
			//g.Log().Errorf(ctx, "failed to run command: %s, run command failed counts:%d", stderr, config.RunCommandFailedCounts)
			return strings.TrimSpace(ret), err
		}
		return s.RunCommand(config)
	}

	return strings.TrimSpace(ret), nil
}

// GetSystemInfo 获取系统信息
func (s *Controller) GetSystemInfo(ctx context.Context) (systemInfo *consts.SystemInfo, err error) {
	// 目前只支持Linux和Mac系统
	var (
		command string
		out     string
	)
	command = "uname"
	out, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return systemInfo, err
	}
	// out = strings.TrimSpace(out)
	switch out {
	case "Linux":
		return s.GetLinuxSystemInfo(ctx)
	case "Darwin":
		return s.GetDarwinSystemInfo(ctx)
	default:
		return systemInfo, fmt.Errorf("unsupported system: %s", out)
	}
}

// CheckProgram 检查程序是否存在
func (s *Controller) CheckProgram(program string) (bool, error) {
	command := "which " + program
	out, err := s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	return out != "", err
}

// GetLinuxSystemInfo 获取Linux系统信息
func (s *Controller) GetLinuxSystemInfo(ctx context.Context) (systemInfo *consts.SystemInfo, err error) {
	systemInfo = &consts.SystemInfo{}
	var (
		command string
		out     string
	)
	command = "cat /etc/os-release"
	out, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return systemInfo, err
	}
	os, err := utils.ParseToJSON(out, "=")
	if err != nil {
		return systemInfo, err
	}
	err = json.Unmarshal([]byte(os), &systemInfo)
	if err != nil {
		return systemInfo, err
	}
	command = "uname -m"
	out, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return systemInfo, err
	}
	systemInfo.Architecture = out
	command = "uname -r"
	out, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return systemInfo, err
	}
	systemInfo.KernelVersion = out
	switch systemInfo.ID {
	case "centos", "bclinux":
		systemInfo.OSVersion, err = s.getCentosSystemDetailVersion(*systemInfo)
		if err != nil {
			return systemInfo, err
		}
	case "ubuntu", "debian":
		systemInfo.OSVersion, err = s.getDebianSystemDetailVersion(*systemInfo)
		if err != nil {
			return systemInfo, err
		}
	default:
		systemInfo.OSVersion = systemInfo.VersionID
	}
	return systemInfo, nil
}

// GetDarwinSystemInfo 获取Mac系统信息
func (s *Controller) GetDarwinSystemInfo(ctx context.Context) (systemInfo *consts.SystemInfo, err error) {
	systemInfo = &consts.SystemInfo{}
	var (
		command string
		out     string
	)
	command = "sw_vers"
	out, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return systemInfo, err
	}
	type DarwinOS struct {
		ProductName    string `json:"product_name"`
		ProductVersion string `json:"product_version"`
		BuildVersion   string `json:"build_version"`
	}
	var darwinOS DarwinOS
	os, err := utils.ParseToJSON(out, ":")
	if err != nil {
		return systemInfo, err
	}
	err = json.Unmarshal([]byte(os), &darwinOS)
	if err != nil {
		return systemInfo, err
	}
	systemInfo.Name = darwinOS.ProductName
	systemInfo.Version = darwinOS.ProductVersion
	systemInfo.ID = "darwin"
	systemInfo.VersionID = darwinOS.ProductVersion
	systemInfo.PrettyName = fmt.Sprintf("%s %s", darwinOS.ProductName, darwinOS.ProductVersion)
	systemInfo.HomeUrl = "https://www.apple.com/macos/"
	systemInfo.BuildID = darwinOS.BuildVersion
	command = "uname -m"
	out, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return systemInfo, err
	}
	systemInfo.Architecture = out
	command = "uname -r"
	out, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return systemInfo, err
	}
	systemInfo.KernelVersion = out
	systemInfo.OSVersion = systemInfo.VersionID
	return systemInfo, nil
}

// GetInstalledPackageList 获取已经安装的软件包列表
func (s *Controller) GetInstalledPackageList(ctx context.Context, in consts.GetInstalledPackageListParam) (out consts.GetInstalledPackageListResult, err error) {
	if s.systemInfo == nil {
		s.systemInfo, err = s.GetSystemInfo(ctx)
		if err != nil {
			return out, err
		}

	}
	switch s.systemInfo.ID {
	case "centos", "bclinux":
		return s.getCentosInstalledPackageList(ctx, in)
	case "ubuntu", "debian":
		return s.getDebianInstalledPackageList(ctx, in)
	default:
		return out, fmt.Errorf("unsupported system: %s", s.systemInfo.ID)
	}
}

// GetPackageInfo 获取某个软件包的详细信息
func (s *Controller) GetPackageInfo(ctx context.Context, in consts.GetPackageInfoParam) (out consts.GetPackageInfoResult, err error) {
	if s.systemInfo == nil {
		s.systemInfo, err = s.GetSystemInfo(ctx)
		if err != nil {
			return out, err
		}
	}
	switch s.systemInfo.ID {
	case "centos", "bclinux":
		return s.getCentosPackageInfo(ctx, in)
	case "ubuntu", "debian":
		return s.getDebianPackageInfo(ctx, in)
	default:
		return out, fmt.Errorf("unsupported system: %s", s.systemInfo.ID)
	}
}

// InstallPackages 安装软件包列表
func (s *Controller) InstallPackages(ctx context.Context, in consts.InstallPackageParam) (installPackageResult consts.InstallPackageResult, err error) {
	// packages := strings.Join(packageList, " ")
	switch s.systemInfo.ID {
	case "centos", "bclinux":
		return s.installCentosPackage(ctx, in)
	// case "ubuntu":
	// return s.DebianInstallPackage(ctx, packageList)
	default:
		return installPackageResult, fmt.Errorf("unsupported system: %s", s.systemInfo.ID)
	}
	// return successfullyInstallPackageList, failedInstallPackageList, nil
	// return installPackageResult, nil
}

// GetPackageFileList 获取某个软件包的文件列表
func (s *Controller) GetPackageFileList(ctx context.Context, in consts.GetPackageFileListParam) (out consts.GetPackageFileListResult, err error) {
	switch s.systemInfo.ID {
	case "centos", "bclinux":
		return s.getCentosPackageFileList(ctx, in)
	case "ubuntu", "debian":
		return s.getDebianPackageFileList(ctx, in)
	default:
		return out, fmt.Errorf("unsupported system: %s", s.systemInfo.ID)
	}
}

// GetPackagesFileList 获取多个软件包文件列表
func (s *Controller) GetPackagesFileList(ctx context.Context, in consts.GetPackagesFileListParam) (out consts.GetPackagesFileListResult, err error) {
	if in.Async {
		return s.GetPackagesFileListAsync(ctx, in)
	}
	for _, packageName := range in.PackageList {
		packageFileList, err := s.GetPackageFileList(ctx, consts.GetPackageFileListParam{
			PackageName: packageName,
		})
		if err != nil {
			return out, err
		}
		out.PackageList = append(out.PackageList, packageFileList)
	}
	return out, nil
}

// GetPackagesFileListAsync 异步获取多个软件包文件列表
func (s *Controller) GetPackagesFileListAsync(ctx context.Context, in consts.GetPackagesFileListParam) (out consts.GetPackagesFileListResult, err error) {
	items := utils.ConvertSliceToInterfaceSlice(in.PackageList)

	AsyncOut, err := utils.AsyncCall(ctx, consts.AsyncCallParam{
		Operation: func(ctx context.Context, item interface{}) (interface{}, error) {
			return s.GetPackageFileList(ctx, consts.GetPackageFileListParam{
				PackageName: item.(string),
			})
		},
		Items: items,
	})
	if err != nil {
		return out, err
	}
	for _, v := range AsyncOut.RetList {
		// g.Log().Debug(ctx, "AsyncOut.RetList:", v)
		out.PackageList = append(out.PackageList, v.Ret.(consts.GetPackageFileListResult))
	}
	return out, nil
}
