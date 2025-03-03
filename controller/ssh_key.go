package controller

import (
	"context"
	"fmt"

	"github.com/ChamHerry/oshelper/consts"
	"github.com/gogf/gf/v2/frame/g"
)

// GetSSHPublicKey 获取ssh公钥
func (s *Controller) GetSSHPublicKey(ctx context.Context, in consts.GetSSHKeyParam) (out consts.GetSSHKeyResult, err error) {
	filePath := getSSHKeyFilePath(in.KeyType)
	if filePath == "" {
		return out, fmt.Errorf("unsupported key type: %s", in.KeyType)
	}
	if !s.IsFileExist(ctx, filePath) {
		if in.Generate {
			return s.GenerateSSHKey(ctx, in)
		}
		return out, fmt.Errorf("file not found: %s", filePath)
	}
	out.FilePath = filePath
	pubKeyFilePath := filePath + ".pub"
	pubKeyFileContent, err := s.GetFileContent(ctx, pubKeyFilePath)
	if err != nil {
		return out, err
	}
	out.PubKeyContent = pubKeyFileContent
	return out, nil
}

// GenerateSSHKey 生成ssh密钥
func (s *Controller) GenerateSSHKey(ctx context.Context, in consts.GetSSHKeyParam) (out consts.GetSSHKeyResult, err error) {
	filePath := getSSHKeyFilePath(in.KeyType)
	if filePath == "" {
		return out, fmt.Errorf("unsupported key type: %s", in.KeyType)
	}
	if s.IsFileExist(ctx, filePath) {
		pubKeyFilePath := filePath + ".pub"
		pubKeyFileContent, err := s.GetFileContent(ctx, pubKeyFilePath)
		if err != nil {
			return out, err
		}
		out.PubKeyContent = pubKeyFileContent
		out.FilePath = filePath
		return out, nil
	}
	// 生成ssh密钥
	if in.Comment == "" {
		in.Comment = "DefaultCommentCreatedByOSHelper"
	}
	command := "ssh-keygen -t " + in.KeyType + " -C " + in.Comment + " -f " + filePath + " -N ''"
	_, err = s.RunCommand(consts.RunCommandConfig{
		Command: command,
	})
	if err != nil {
		return out, err
	}
	return out, nil
}

// getSSHKeyFilePath 获取ssh密钥文件路径
func getSSHKeyFilePath(keyType string) string {
	switch keyType {
	case "dsa":
		return "~/.ssh/id_dsa"
	case "ecdsa":
		return "~/.ssh/id_ecdsa"
	case "ed25519", "":
		return "~/.ssh/id_ed25519"
	case "rsa":
		return "~/.ssh/id_rsa"
	case "rsa1":
		return "~/.ssh/id_rsa1"
	default:
		return ""
	}
}

// 备份AuthorizedKeys文件
func (s *Controller) BackupAuthorizedKeys(ctx context.Context) error {
	checkBackupCmd := "test -f ~/.ssh/authorized_keys.migration_scheduling.bak"
	_, err := s.RunCommand(consts.RunCommandConfig{
		Command: checkBackupCmd,
	})
	if err == nil {
		// 备份文件已存在，无需重复备份
		return nil
	}
	g.Log().Debug(ctx, "backup authorized keys")
	return nil
	// authorizedKeysFilePath := "~/.ssh/authorized_keys"
	// if !s.IsFileExist(ctx, authorizedKeysFilePath) {
	// 	return fmt.Errorf("file not found: %s", authorizedKeysFilePath)
	// }
	// backupFilePath := authorizedKeysFilePath + ".backup"
	// if s.IsFileExist(ctx, backupFilePath) {
	// 	return fmt.Errorf("backup file already exists: %s", backupFilePath)
	// }
	// return s.CopyFile(ctx, authorizedKeysFilePath, backupFilePath)
}

// 恢复AuthorizedKeys文件
func (s *Controller) RestoreAuthorizedKeys(ctx context.Context) error {
	// authorizedKeysFilePath := "~/.ssh/authorized_keys"
	// backupFilePath := authorizedKeysFilePath + ".backup"
	// if !s.IsFileExist(ctx, backupFilePath) {
	// 	return fmt.Errorf("backup file not found: %s", backupFilePath)
	// }
	// return s.CopyFile(ctx, backupFilePath, authorizedKeysFilePath)
	return nil
}
