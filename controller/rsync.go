package controller

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ChamHerry/oshelper/consts"

	"github.com/gogf/gf/v2/frame/g"
)

// RsyncLocalToRemote 将本地文件同步到远程服务器
func (s *Controller) RsyncLocalToRemote(ctx context.Context, sshConfig consts.SSHConfig, localPath, remotePath string) error {
	rsync, err := s.CheckProgram("rsync")
	g.Log().Debugf(ctx, "rsync:%v, err:%v", rsync, err)
	if err != nil {
		return err
	}
	if !rsync {
		return fmt.Errorf("rsync not found")
	}
	remoteHost := sshConfig.Username + "@" + sshConfig.IPAddress
	command := "rsync -ar --inplace -e 'ssh -p " + strconv.Itoa(sshConfig.Port) + " -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no' " + localPath + " " + remoteHost + ":" + remotePath
	g.Log().Debugf(ctx, "rsync local path: %s, remote path: %s", localPath, remotePath)
	_, err = s.RunCommand(consts.RunCommandConfig{
		Command: command,
	})
	if err != nil {
		errLines := strings.Split(err.Error(), "\n")
		if (strings.Contains(err.Error(), "Warning") || strings.Contains(err.Error(), "警告")) && len(errLines) == 1 {
			return nil
		}
		err = errors.New(strings.Join(errLines[1:], "\n"))
		g.Log().Warning(ctx, "RsyncLocalToRemote error:", err)
		return err
	}
	return nil
}

// RsyncRemoteToLocal 将远程服务器文件同步到本地
func (s *Controller) RsyncRemoteToLocal(ctx context.Context, sshConfig consts.SSHConfig, remotePath, localPath string) error {
	rsync, err := s.CheckProgram("rsync")
	g.Log().Debugf(ctx, "rsync:%v, err:%v", rsync, err)
	if err != nil {
		return err
	}
	if !rsync {
		return fmt.Errorf("rsync not found")
	}
	remoteHost := sshConfig.Username + "@" + sshConfig.IPAddress
	command := "rsync -ar --inplace -e 'ssh -p " + strconv.Itoa(sshConfig.Port) + " -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no' " + remoteHost + ":" + remotePath + " " + localPath
	g.Log().Debugf(ctx, "rsync remote path: %s, local path: %s", remotePath, localPath)
	_, err = s.RunCommand(consts.RunCommandConfig{
		Command: command,
	})
	if err != nil {
		errLines := strings.Split(err.Error(), "\n")
		if (strings.Contains(err.Error(), "Warning") || strings.Contains(err.Error(), "警告")) && len(errLines) == 1 {
			return nil
		}
		err = errors.New(strings.Join(errLines[1:], "\n"))
		g.Log().Warning(ctx, "RsyncLocalToRemote error:", err)
		return err
	}
	return nil
}
