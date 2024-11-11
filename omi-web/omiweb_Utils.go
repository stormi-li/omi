package omiweb

import (
	"io/fs"
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

func copyEmbeddedFiles() error {
	srcFS := templateSource
	destDir := "static"
	// 遍历嵌入文件系统中的所有文件
	err := fs.WalkDir(srcFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录，仅复制文件
		if d.IsDir() {
			return nil
		}

		// 确保目标文件夹路径
		destPath := filepath.Join(destDir, filepath.Base(path))
		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		// 检查目标文件是否已存在
		if _, err := os.Stat(destPath); err == nil {
			return nil
		} else if !os.IsNotExist(err) {
			return err
		}

		// 读取嵌入文件内容
		content, err := srcFS.ReadFile(path)
		if err != nil {
			return err
		}

		// 写入文件到目标文件夹
		if err := os.WriteFile(destPath, content, os.ModePerm); err != nil {
			return err
		}

		return nil
	})

	return err
}
