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
