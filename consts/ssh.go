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
