package consts

type DownloadFileParam struct {
	URL      string
	DirPath  string
	FileName string
}

type DownloadFileResult struct {
	FilePath string
}

type DeleteFileParam struct {
	FilePath string
}

type DeleteFileResult struct {
}
