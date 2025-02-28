package consts

// PackageInfo 包信息
type PackageInfo struct {
	FullName      string   `json:"full_name,omitempty"`
	Name          string   `json:"name,omitempty"`
	Version       string   `json:"version,omitempty"`
	ReleaseNumber string   `json:"release_number,omitempty"`
	OS            string   `json:"os,omitempty"`
	Architecture  string   `json:"architecture,omitempty"`
	Description   string   `json:"description,omitempty"`
	Files         []string `json:"files,omitempty"`
	Requires      []string `json:"requires,omitempty"`
	Provides      []string `json:"provides,omitempty"`
}

// InstallPackageParam 安装软件包入参
type InstallPackageParam struct {
	PackageList []PackageInfo
}

// InstallPackageResult 安装软件包结果
type InstallPackageResult struct {
	IgnoredInstallPackageList      []PackageInfo `json:"ignored_install_package_list"`
	SuccessfullyInstallPackageList []PackageInfo `json:"successfully_install_package_list"`
	FailedInstallPackageList       []PackageInfo `json:"failed_install_package_list"`
	Total                          int           `json:"total"`
}

// 获取已经安装的软件包列表入参
type GetInstalledPackageListParam struct {
}

// 获取已经安装的软件包列表出参
type GetInstalledPackageListResult struct {
	PackageList []string
}

// 获取某个软件包的文件列表入参
type GetPackageFileListParam struct {
	PackageName string
}

// 获取某个软件包的文件列表出参
type GetPackageFileListResult struct {
	PackageInfo
}

// 获取某个软件包的详细信息入参
type GetPackageInfoParam struct {
	PackageName string
}

// 获取某个软件包的详细信息出参
type GetPackageInfoResult struct {
	PackageInfo
}

// 获取多个软件包文件列表入参
type GetPackagesFileListParam struct {
	PackageList []string
	Async       bool
}

// 获取多个软件包文件列表出参
type GetPackagesFileListResult struct {
	PackageList []GetPackageFileListResult
}
