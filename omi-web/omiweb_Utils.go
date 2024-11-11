package omiweb

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func copyResource(src string) {
	err := copyDir(src, target_path)
	if err != nil {
		log.Fatalln(err)
	}
}

func getSourceFilePath() string {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	return dir
}

func copyFile(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		// 如果目标文件已存在，不进行复制
		return nil
	} else if !os.IsNotExist(err) {
		// 如果发生其他错误，返回错误
		return err
	}

	// 读取源文件内容
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// 写入目标文件
	err = os.WriteFile(dst, data, 0644) // 使用 0644 权限创建目标文件
	return err
}

// copyDir 递归复制整个文件夹
func copyDir(srcDir, dstDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 构建目标路径
		relativePath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dstDir, relativePath)

		if info.IsDir() {
			// 创建目标文件夹
			return os.MkdirAll(dstPath, info.Mode())
		} else {
			// 复制文件
			return copyFile(path, dstPath)
		}
	})
}
