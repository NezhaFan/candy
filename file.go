package util

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

// 判断文件是否存在
func IsFileExist(name string) bool {
	_, err := os.Stat(name)
	return err == nil || os.IsExist(err)
}

// 对文件追加写 (文件不存在则创建)
func AppendFile(filename string) (*os.File, error) {
	dir := filepath.Dir(filename)
	if !IsFileExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)
	return f, err
}

// 清除文件内容 (文件不存在则创建)
func TruncateFile(name string) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	return file.Truncate(0)
}

// 按行读取文件
func ReadFileByLine(name string) ([]string, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	strs := make([]string, 0, 256)
	r := bufio.NewReader(file)
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return strs, err
		}
		strs = append(strs, s)
	}
}
