package consts

// SSHConfig 用于存储 SSH 配置
type SSHConfig struct {
	Username  string `json:"username",omitempty,default:"root"`
	Password  string `json:"password",omitempty,default:""`
	IPAddress string `json:"ip_address",omitempty,default:""`
	Port      int    `json:"port",omitempty,default:22`
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
	Generate bool   `json:"generate",omitempty,default:false` // 是否生成公钥
	KeyType  string `json:"key_type",omitempty`               // 公钥类型，可选值为 dsa | ecdsa | ed25519 | rsa | rsa1
	Comment  string `json:"comment",omitempty`                // 公钥备注

}

// GetSSHKeyResult 获取ssh公钥结果
type GetSSHKeyResult struct {
	FilePath      string `json:"file_path",omitempty`       // 公钥文件路径
	PubKeyContent string `json:"pub_key_content",omitempty` // 公钥内容
}
