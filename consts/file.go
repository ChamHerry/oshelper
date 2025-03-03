package consts

type DownloadFileParam struct {
	URL      string
	DirPath  string
	FileName string
}

// DownloadFileResult 下载文件结果
type DownloadFileResult struct {
	FilePath string
}

// DeleteFileParam 删除文件参数
type DeleteFileParam struct {
	FilePath string
}

// DeleteFileResult 删除文件结果
type DeleteFileResult struct {
}

// WriteFileParam 写入文件参数
type WriteFileParam struct {
	FilePath  string
	Content   string
	Overwrite bool
}
