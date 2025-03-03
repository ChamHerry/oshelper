package controller

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ChamHerry/oshelper/consts"
	"github.com/ChamHerry/oshelper/utils"
	"github.com/gogf/gf/v2/frame/g"
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

// CopyFile 复制文件
func (s *Controller) CopyFile(ctx context.Context, in, out string) (err error) {
	g.Log().Infof(ctx, "CopyFile: %s -> %s", in, out)
	if !s.IsFileOrDirExist(ctx, in) {
		return fmt.Errorf("file not found: %s", in)
	}
	if s.IsFileOrDirExist(ctx, out) {
		return fmt.Errorf("file already exists: %s", out)
	}
	command := "cp -r " + in + " " + out
	_, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	return err
}

// CreateDir 创建文件夹
func (s *Controller) CreateDir(ctx context.Context, in string) (err error) {
	if s.IsFileOrDirExist(ctx, in) {
		return nil
	}
	command := "mkdir -p " + in
	_, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	return err
}

// CreateFile 创建文件
func (s *Controller) CreateFile(ctx context.Context, in string) (err error) {
	var command string
	if !s.IsFileOrDirExist(ctx, filepath.Dir(in)) {
		err = s.CreateDir(ctx, filepath.Dir(in))
		if err != nil {
			return err
		}
	}
	command = "touch " + in
	_, err = s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	return err
}

// WriteFile 写入文件
func (s *Controller) WriteFile(ctx context.Context, in consts.WriteFileParam) (err error) {
	if !s.IsFileOrDirExist(ctx, in.FilePath) {
		err = s.CreateFile(ctx, in.FilePath)
		if err != nil {
			return err
		}
	}
	if in.Overwrite {
		command := "echo \"" + in.Content + "\" > " + in.FilePath
		_, err = s.RunCommand(consts.RunCommandConfig{
			Command:                command,
			RunCommandFailedCounts: 0,
		})
	} else {
		command := "echo \"" + in.Content + "\" >> " + in.FilePath
		_, err = s.RunCommand(consts.RunCommandConfig{
			Command:                command,
			RunCommandFailedCounts: 0,
		})
	}
	return err
}
