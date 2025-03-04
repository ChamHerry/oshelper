package consts

const (
	SSHKeyTypeED25519               = "ed25519"
	SSHKeyTypeRSA                   = "rsa"
	SSHKeyTypeECDSA                 = "ecdsa"
	SSHKeyTypeDSA                   = "dsa"
	SSHKeyTypeRSA1                  = "rsa1"
	SSHBackupAuthorizedKeysFilePath = "~/.ssh/authorized_keys.migration_scheduling.bak"
	SSHAuthorizedKeysFilePath       = "~/.ssh/authorized_keys"
)

// SSHConfig 用于存储 SSH 配置
type SSHConfig struct {
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	IPAddress string `json:"ip_address,omitempty"`
	Port      int    `json:"port,omitempty"`
}

// DefaultSSHConfig 默认SSH配置
var DefaultSSHConfig = SSHConfig{
	Username:  "root",
	Password:  "",
	IPAddress: "",
	Port:      22,
}

// GetSSHKeyParam 获取ssh公钥入参
type GetSSHKeyParam struct {
	Generate bool   `json:"generate,omitempty"` // 是否生成公钥
	KeyType  string `json:"key_type,omitempty"` // 公钥类型，可选值为 dsa | ecdsa | ed25519 | rsa | rsa1
	Comment  string `json:"comment,omitempty"`  // 公钥备注
}

// GetSSHKeyResult 获取ssh公钥结果
type GetSSHKeyResult struct {
	FilePath      string `json:"file_path,omitempty"`       // 公钥文件路径
	PubKeyContent string `json:"pub_key_content,omitempty"` // 公钥内容
}

// AddSSHKeyParam 添加免密认证参数
type AddSSHKeyParam struct {
	SSHConfig SSHConfig `json:"ssh_config,omitempty"` // ssh配置
	KeyType   string    `json:"key_type,omitempty"`   // 公钥类型，可选值为 dsa | ecdsa | ed25519 | rsa | rsa1
	Comment   string    `json:"comment,omitempty"`    // 公钥备注
	IsBackup  bool      `json:"is_backup,omitempty"`  // 是否备份authorized_keys文件
}

// DeleteSSHKeyParam 删除免密认证参数
type DeleteSSHKeyParam struct {
	SSHConfig SSHConfig `json:"ssh_config,omitempty"` // ssh配置
}
