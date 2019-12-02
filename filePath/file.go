package filePath

import (
	"path/filepath"
	"strings"
)

type FilePath interface {
	GetFileDir() string
	GetFileName() string
	GetFileType() string
	GetFilePath() string
	Format(string) (string, bool)
}
type filePath struct {
	path string
	dir  string
	name string
	typ  string
}

func (f *filePath) GetFileDir() string {
	return f.dir
}

func (f *filePath) GetFileName() string {
	return f.name
}

func (f *filePath) GetFileType() string {
	return f.typ
}

func (f *filePath) GetFilePath() string {
	return f.path
}

func (f *filePath) Format(formatter string) (string, bool) {
	var count int = 0
	// @FILEDIR/@FILENAME.@FILETYPE
	if strings.Contains(formatter, "@FILEDIR") {
		count++
		formatter = strings.ReplaceAll(formatter, "@FILEDIR", f.dir)
	}
	if strings.Contains(formatter, "@FILENAME") {
		count++
		formatter = strings.ReplaceAll(formatter, "@FILENAME", f.name)
	}
	if strings.Contains(formatter, "@FILETYPE") {
		count++
		formatter = strings.ReplaceAll(formatter, "@FILETYPE", f.typ)
	}
	if count == 3 {
		return f.path, true
	}
	return formatter, false
}

func NewFilePath(path string) FilePath {
	return &filePath{
		path: path,
		dir:  filepath.Dir(path),
		name: strings.ReplaceAll(filepath.Base(path), filepath.Ext(path), ""),
		typ:  filepath.Ext(path),
	}
}
