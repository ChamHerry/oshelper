package consts

// SystemInfo 用于存储系统基本信息
type SystemInfo struct {
	FullName        string `json:"full_name,omitempty"`
	Name            string `json:"name,omitempty"`
	Version         string `json:"version,omitempty"`
	ID              string `json:"id,omitempty"`
	IDLike          string `json:"id_like,omitempty"`
	VersionCodename string `json:"version_codename,omitempty"`
	VersionID       string `json:"version_id,omitempty"`
	PrettyName      string `json:"pretty_name,omitempty"`
	AnsiColor       string `json:"ansi_color,omitempty"`
	CPEName         string `json:"cpe_name,omitempty"`
	HomeUrl         string `json:"home_url,omitempty"`
	BuildID         string `json:"build_id,omitempty"`
	Variant         string `json:"variant,omitempty"`
	VariantID       string `json:"variant_id,omitempty"`
	OSVersion       string `json:"os_version,omitempty"`
	Architecture    string `json:"architecture,omitempty"`
	KernelVersion   string `json:"kernel_version,omitempty"`
}

// 文件信息
type FileInfo struct {
	Name         string `json:"name,omitempty"`         // 文件名
	Size         int    `json:"size,omitempty"`         // 文件大小
	Mode         string `json:"mode,omitempty"`         // 文件权限
	ModTime      string `json:"mod_time,omitempty"`     // 文件修改时间
	FileType     string `json:"file_type,omitempty"`    // 文件类型
	User         string `json:"user,omitempty"`         // 文件所属用户
	Group        string `json:"group,omitempty"`        // 文件所属组
	LinkPath     string `json:"link_path,omitempty"`    // 链接路径
	GlobalPath   string `json:"global_path,omitempty"`  // 全局路径
	Architecture string `json:"architecture,omitempty"` // 架构
}

// GetSystemFileParam 获取系统文件参数
type GetSystemFileParam struct {
	FilePath string `json:"file_path,omitempty"` // 文件路径
	Async    bool   `json:"async,omitempty"`     // 是否异步
}

// GetSystemFileResult 获取系统文件结果
type GetSystemFileResult struct {
	FilePath string `json:"file_path,omitempty"` // 文件路径
	Content  string `json:"content,omitempty"`   // 文件内容
}

// GetFilePathListParam 获取某个路径下的文件清单参数
type GetFilePathListParam struct {
	FilePath string `json:"file_path,omitempty"` // 文件路径
}

// GetFilePathListResult 获取某个路径下的文件清单结果
type GetFilePathListResult struct {
	FileInfoList []FileInfo `json:"file_info_list,omitempty"` // 文件信息列表
}
