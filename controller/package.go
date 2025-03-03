package controller

import (
	"context"
	"fmt"

	"github.com/ChamHerry/oshelper/consts"
	"github.com/ChamHerry/oshelper/utils"
)

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

// InstallPackages 安装软件包列表
func (s *Controller) InstallPackages(ctx context.Context, in consts.InstallPackagesParam) (installPackageResult consts.InstallPackagesResult, err error) {
	// packages := strings.Join(packageList, " ")
	switch s.systemInfo.ID {
	case "centos", "bclinux":
		return s.installCentosPackages(ctx, in)
	// case "ubuntu":
	// return s.DebianInstallPackage(ctx, packageList)
	default:
		return installPackageResult, fmt.Errorf("unsupported system: %s", s.systemInfo.ID)
	}
	// return successfullyInstallPackageList, failedInstallPackageList, nil
	// return installPackageResult, nil
}
