package controller

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
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

// FindFile 查找文件
func (s *Controller) FindFile(ctx context.Context, in consts.FindFileParam) (out consts.FindFileResult, err error) {
	command := "find " + in.DirPath + " -name \"" + in.FileName + "\" 2>/dev/null"
	content, err := s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return out, err
	}
	out.FilePathList = strings.Split(strings.TrimSpace(content), "\n")
	return out, nil
}

// GetFileInfo 获取文件信息
func (s *Controller) GetFileInfo(ctx context.Context, in consts.GetFileInfoParam) (out consts.GetFileInfoResult, err error) {
	basePath := filepath.Dir(in.FilePath)
	fileName := filepath.Base(in.FilePath)
	command := "ls -l " + strings.TrimRight(basePath, "/") + "/ | grep -w \"" + fileName + "\""
	g.Log().Debug(ctx, "command:", command)
	commandOut, err := s.RunCommand(consts.RunCommandConfig{
		Command:                command,
		RunCommandFailedCounts: 0,
	})
	if err != nil {
		return out, err
	}
	lines := strings.Split(commandOut, "\n")
	if len(lines) == 0 {
		return out, fmt.Errorf("file not found: %s", in.FilePath)
	}
	fields := strings.Fields(lines[0])
	if len(fields) < 9 {
		return out, fmt.Errorf("file info format error: %s", lines[0])
	}
	out.FileInfo.Name = fields[8]
	out.FileInfo.Size, err = strconv.Atoi(fields[4])
	if err != nil {
		return out, err
	}
	out.FileInfo.Mode = fields[0][:10]
	out.FileInfo.ModTime = fields[5] + " " + fields[6]
	switch fields[0][0] {
	case 'd':
		out.FileInfo.FileType = "directory"
		out.Architecture = "noarch"
	case '-':
		out.FileInfo.FileType = "file"
		// 获取架构
		command = "file " + in.FilePath
		g.Log().Debug(ctx, "command:", command)
		commandOut, err = s.RunCommand(consts.RunCommandConfig{
			Command:                command,
			RunCommandFailedCounts: 0,
		})
		if err != nil {
			return out, err
		}
		commandOut = strings.TrimSpace(commandOut)
		if strings.Contains(commandOut, "x86-64") || strings.Contains(commandOut, "x86_64") {
			out.FileInfo.Architecture = "x86_64"
		} else if strings.Contains(commandOut, "ARM") || strings.Contains(commandOut, "aarch64") {
			out.FileInfo.Architecture = "aarch64"
		} else {
			out.FileInfo.Architecture = "noarch"
		}
	case 'l':
		out.FileInfo.FileType = "link"
		out.FileInfo.LinkPath, err = s.GetRealPath(ctx, in.FilePath)
		if err != nil {
			return out, err
		}
		out.Architecture = "noarch"
	}
	out.FileInfo.User = fields[2]
	out.FileInfo.Group = fields[3]
	out.GlobalPath = in.FilePath
	// // 获取架构
	// var archPath string
	// if out.FileInfo.FileType == "link" {
	// 	archPath = out.FileInfo.LinkPath
	// } else {
	// 	archPath = in.FilePath
	// }

	return out, nil
}

// GetFileInfoByStat 获取文件信息
// func (s *Controller) GetFileInfoByStat(ctx context.Context, filePath string) (out consts.GetFileInfoResult, err error) {
// 	command := "stat -c %A,%s,%U,%G,%h,%i,%n,%F,%y,%z,%b,%X,%Y,%Z,%a,%b,%c,%d,%f,%g,%h,%i,%n,%o,%p,%s,%t,%u,%w,%x,%y,%z " + filePath
// 	commandOut, err := s.RunCommand(consts.RunCommandConfig{
// 		Command:                command,
// 		RunCommandFailedCounts: 0,
// 	})
// 	if err != nil {
// 		return out, err
// 	}
// 	lines := strings.Split(commandOut, "\n")
// 	if len(lines) == 0 {
// 		return out, fmt.Errorf("file not found: %s", filePath)
// 	}
// 	fields := strings.Fields(lines[0])
// 	if len(fields) < 30 {
// 		return out, fmt.Errorf("file info format error: %s", lines[0])
// 	}
// 	out.FileInfo.Name = fields[8]
// 	out.FileInfo.Size, err = strconv.Atoi(fields[1])
// 	if err != nil {
// 		return out, err
// 	}
// 	out.FileInfo.Mode = fields[0]
// 	out.FileInfo.ModTime = fields[5] + " " + fields[6]
// 	switch fields[0][0] {
// 	case 'd':
// 		out.FileInfo.FileType = "directory"
// 	case '-':
// 		out.FileInfo.FileType = "file"
// 	case 'l':
// 		out.FileInfo.FileType = "link"
// 		out.FileInfo.LinkPath, err = s.GetRealPath(ctx, filePath)
// 		if err != nil {
// 			return out, err
// 		}
// 	}
// 	out.FileInfo.User = fields[2]
// 	out.FileInfo.Group = fields[3]
// 	out.GlobalPath = filePath
// 	// 获取架构
// 	var archPath string
// 	if out.FileInfo.FileType == "link" {
// 		archPath = out.FileInfo.LinkPath
// 	} else {
// 		archPath = filePath
// 	}
// 	command = "file " + archPath
// 	commandOut, err = s.RunCommand(consts.RunCommandConfig{
// 		Command:                command,
// 		RunCommandFailedCounts: 0,
// 	})
// 	if err != nil {

// }

// 获取软连接的真实路径
func (s *Controller) GetRealPath(ctx context.Context, filePath string) (string, error) {
	commandConfig := consts.RunCommandConfig{
		Command: fmt.Sprintf("readlink -f %s", filePath),
	}
	out, err := s.RunCommand(commandConfig)
	if err != nil {
		return "", fmt.Errorf("failed to get real path: %s", err)
	}
	out = strings.TrimSpace(out)
	return out, nil
}
