package os

import (
	"os"
)

func FileExist(filename string) (bool, error) {

	// 使用os.Stat()函数获取文件的信息
	_, err := os.Stat(filename)

	// 判断文件是否存在
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}
