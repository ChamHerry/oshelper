package controller

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ChamHerry/oshelper/consts"
	"github.com/ChamHerry/oshelper/utils"
)

// DownloadFile 下载文件
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

// DeleteFile 删除文件
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

// IsFileExist 判断是否存在文件
func (s *Controller) IsFileExist(ctx context.Context, in string) (out bool) {
	command := "ls -a " + in
	if _, err := s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	}); err != nil {

		return false
	}
	return true
}

// IsDirExist 判断是否存在文件夹
func (s *Controller) IsDirExist(ctx context.Context, in string) (out bool) {
	in = strings.TrimRight(in, "/")
	command := "ls -a " + filepath.Dir(in) + " | grep \"" + filepath.Base(in) + "\""
	if _, err := s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	}); err != nil {
		return false
	}
	return true
}

// IsFileOrDirExist 判断是否存在文件或文件夹
func (s *Controller) IsFileOrDirExist(ctx context.Context, in string) (out bool) {
	command := "test -e " + in
	if _, err := s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	}); err != nil {
		return false
	}
	return true
}

// GetFileContent 获取文件内容
func (s *Controller) GetFileContent(ctx context.Context, in string) (out string, err error) {
	if !s.IsFileExist(ctx, in) {
		return out, fmt.Errorf("file not found: %s", in)
	}
	command := "cat " + in
	out, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return out, err
	}
	return out, nil
}
