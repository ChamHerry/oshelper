package controller

import (
	"context"
	"fmt"
	"strings"

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

// BackupAuthorizedKeys 备份AuthorizedKeys文件
func (s *Controller) BackupAuthorizedKeys(ctx context.Context) error {
	exit := s.IsFileOrDirExist(ctx, "~/.ssh/authorized_keys.migration_scheduling.bak")
	if exit {
		// 备份文件已存在，无需重复备份
		return nil
	}
	err := s.CopyFile(ctx, "~/.ssh/authorized_keys", "~/.ssh/authorized_keys.migration_scheduling.bak")
	if err != nil {
		return err
	}
	return nil
}

// RestoreAuthorizedKeys 恢复AuthorizedKeys文件
func (s *Controller) RestoreAuthorizedKeys(ctx context.Context) error {
	backupFilePath := "~/.ssh/authorized_keys.migration_scheduling.bak"
	if !s.IsFileOrDirExist(ctx, backupFilePath) {
		return fmt.Errorf("backup file not found: %s", backupFilePath)
	}
	return s.CopyFile(ctx, backupFilePath, "~/.ssh/authorized_keys")
}

// AddSSHKey 添加免密认证
func (s *Controller) AddSSHKey(ctx context.Context, in consts.AddSSHKeyParam) error {
	authorizedKeysFilePath := "~/.ssh/authorized_keys"
	if !s.IsFileOrDirExist(ctx, authorizedKeysFilePath) {
		err := s.CreateFile(ctx, authorizedKeysFilePath)
		if err != nil {
			return err
		}
	}
	sshKeyContent, err := s.GetSSHPublicKey(ctx, consts.GetSSHKeyParam{
		KeyType:  in.KeyType,
		Comment:  in.Comment,
		Generate: true,
	})
	if err != nil {
		return err
	}
	remoteController, err := NewController(in.SSHConfig)
	if err != nil {
		return err
	}
	if in.IsBackup {
		// 备份authorized_keys文件
		g.Log().Debug(ctx, "备份authorized_keys文件")
		err = remoteController.BackupAuthorizedKeys(ctx)
		if err != nil {
			return err
		}
	}
	// 查看authorized_keys文件内容
	authorizedKeysFileContent, err := remoteController.GetFileContent(ctx, authorizedKeysFilePath)
	if err != nil {
		return err
	}
	if strings.Contains(authorizedKeysFileContent, sshKeyContent.PubKeyContent) {
		g.Log().Debug(ctx, "authorized_keys文件已存在，无需重复添加")
		return nil
	}
	// 写入authorized_keys文件
	err = remoteController.WriteFile(ctx, consts.WriteFileParam{
		FilePath:  authorizedKeysFilePath,
		Content:   sshKeyContent.PubKeyContent,
		Overwrite: false,
	})
	if err != nil {
		return err
	}
	return nil
}
