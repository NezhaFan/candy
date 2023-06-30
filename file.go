package candy

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

// 打开文件 (文件不存在则创建)
func OpenFile(filename string) (*os.File, error) {
	dir := filepath.Dir(filename)
	if !IsFileExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return os.OpenFile(filename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0775)
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
			break
		}
		strs = append(strs, s)
	}

	return strs, err
}
