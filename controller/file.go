package controller

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ChamHerry/oshelper/consts"
	"github.com/ChamHerry/oshelper/utils"
)

// 下载文件
func (s *Controller) DownloadFile(ctx context.Context, in consts.DownloadFileParam) (out consts.DownloadFileResult, err error) {
	if in.URL == "" {
		return out, fmt.Errorf("url is empty")
	}
	if in.DirPath == "" {
		// 生成一个随机文件夹
		in.DirPath = "/tmp/download_file_" + utils.RandomString(10)
	}
	// 创建文件夹
	command := "mkdir -p " + in.DirPath
	_, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return out, err
	}
	// 下载文件
	if in.FileName == "" {
		splitList := strings.Split(in.URL, "/")
		in.FileName = splitList[len(splitList)-1]
	}
	command = "wget -O " + filepath.Join(in.DirPath, in.FileName) + " " + in.URL
	_, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return out, err
	}
	out.FilePath = filepath.Join(in.DirPath, in.FileName)
	return out, nil
}

// 删除文件
func (s *Controller) DeleteFile(ctx context.Context, in consts.DeleteFileParam) (out consts.DeleteFileResult, err error) {
	command := "rm -rf " + in.FilePath
	_, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return out, err
	}
	return out, nil
}
