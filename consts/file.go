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

// FindFileParam 查找文件参数
type FindFileParam struct {
	DirPath  string
	FileName string
}

// FindFileResult 查找文件结果
type FindFileResult struct {
	FilePathList []string `json:"file_path_list,omitempty"` // 文件路径列表
}
