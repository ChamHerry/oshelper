package controller

import (
	"context"
	"fmt"
	"oshelper/consts"

	"github.com/gogf/gf/v2/frame/g"
)

func (s *Controller) RsyncLocalToRemote(ctx context.Context, sshConfig consts.SSHConfig, localPath, remotePath string) error {
	rsync, err := s.CheckProgram("rsync1")
	g.Log().Debugf(ctx, "rsync:%v, err:%v", rsync, err)
	if err != nil {
		return err
	}
	if !rsync {
		return fmt.Errorf("rsync not found")
	}
	// controller, err := controller.NewController(sshConfig)
	// if err != nil {
	// 	return err
	// }
	// s.client.
	return nil

}
