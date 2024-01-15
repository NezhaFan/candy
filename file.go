package candy

import (
	"bufio"
	"io"
	"os"
	"path/filepath"
)

// 判断文件是否存在
func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// 打开文件 (文件不存在则创建)
func FileOpen(filename string) (*os.File, error) {
	dir := filepath.Dir(filename)
	if !FileExists(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0775)
}

// 按行读取文件
func ReadByLine(filename string, fn func(string) error) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	r := bufio.NewReader(file)
	for {
		s, err := r.ReadString('\n')
		if err == io.EOF {
			return nil
		}
		if err == nil {
			err = fn(s)
		}
		if err != nil {
			return err
		}
	}
}
