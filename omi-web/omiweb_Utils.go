package omiweb

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 获取当前时间并转换为 UTC 字符串
func getCurrentTimeString() string {
	currentTime := time.Now().UTC() // 设置为 UTC
	return currentTime.Format("2006-01-02 15:04:05")
}

// 将时间字符串解析为 UTC 时间
func parseTimeString(timeString string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	parsedTime, err := time.Parse(layout, timeString)
	if err != nil {
		return time.Time{}, err
	}
	return parsedTime.UTC(), nil // 设置为 UTC
}

// 比较字符串时间和当前时间，判断是否超过 2 秒
func isMoreThanTwoSecondsAgo(timeString string) bool {
	parsedTime, err := parseTimeString(timeString)
	if err != nil {
		return true // 如果解析出错，直接返回 true
	}

	currentTime := time.Now().UTC() // 统一设置为 UTC
	twoSecondsLater := parsedTime.Add(2 * time.Second)

	return currentTime.After(twoSecondsLater)
}

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

func getStringAfterSecondSlash(input string) string {
	// 找到第一个 "/" 的位置
	firstSlashIndex := strings.Index(input, "/")
	if firstSlashIndex == -1 {
		return ""
	}

	// 在第一个 "/" 之后的字符串中再查找第二个 "/"
	secondSlashIndex := strings.Index(input[firstSlashIndex+1:], "/")
	if secondSlashIndex == -1 {
		return ""
	}

	// 计算第二个 "/" 的实际位置并截取子字符串
	secondSlashIndex += firstSlashIndex + 1
	return input[secondSlashIndex+1:]
}

// copyFile 复制单个文件
func copyFile(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		// 如果目标文件已存在，不进行复制
		return nil
	} else if !os.IsNotExist(err) {
		// 如果发生其他错误，返回错误
		return err
	}

	sourceFile, err := os.Open(src)

	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// 使用 io.Copy 复制文件内容
	_, err = io.Copy(destinationFile, sourceFile)
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
