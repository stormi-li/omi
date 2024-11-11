package omiweb

import (
	"container/list"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type cacheItem struct {
	filename string // URL 路径
	filepath string // 文件在磁盘上的路径
	size     int    // 文件大小
}

type fileCache struct {
	lock         sync.RWMutex
	curSize      int
	maxCacheSize int
	cacheDir     string
	itemMap      map[string]*list.Element
	itemList     *list.List // 用于 LRU 策略
}

// 初始化文件缓存系统，读取现有缓存目录内容
func newFileCache(cacheDir string, maxSize int) (*fileCache, error) {
	fileCache := &fileCache{
		itemMap:  make(map[string]*list.Element),
		itemList: list.New(),
	}

	// 设置缓存目录和最大容量
	fileCache.cacheDir = cacheDir
	fileCache.maxCacheSize = maxSize

	// 创建缓存目录（如果不存在）
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	// 遍历缓存目录并加载文件信息
	err := filepath.WalkDir(cacheDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 忽略目录本身
		if d.IsDir() {
			return nil
		}

		// 获取文件信息并构建缓存项
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}

		// 文件大小检查
		fileSize := int(fileInfo.Size())
		if fileSize > maxSize {
			return nil // 跳过超出容量的文件
		}

		filename := filepath.Base(path)
		if log_cache {
			log.Println("新增缓存", filename)
		}
		cacheItem := &cacheItem{filename: filename, filepath: path, size: fileSize}
		elem := fileCache.itemList.PushFront(cacheItem)
		fileCache.itemMap[filename] = elem
		fileCache.curSize += fileSize
		// 超出容量时，清理最旧的缓存文件
		for fileCache.curSize > maxSize {
			fileCache.removeOldest()
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return fileCache, nil
}

// 将文件添加到磁盘缓存，超过容量时使用 LRU 策略清理
func (fileCache *fileCache) UpdateCache(url *url.URL, data []byte) {
	fileCache.lock.Lock()
	defer fileCache.lock.Unlock()
	if len(data) == 0 {
		return
	}
	filename := strings.ReplaceAll(url.String(), "/", "@")

	cachePath := filepath.Join(fileCache.cacheDir, filename)

	// 文件大小超过缓存最大容量时直接返回
	fileSize := len(data)
	if fileSize > fileCache.maxCacheSize {
		return
	}

	// 清理超出容量的旧文件，确保有足够空间
	for fileCache.curSize+fileSize > fileCache.maxCacheSize {
		fileCache.removeOldest()
	}

	// 写入文件到磁盘
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return
	}

	if log_cache {
		log.Println("新增缓存", filename)
	}

	// 新建缓存项并添加到 LRU 列表
	cacheItem := &cacheItem{filename: filename, filepath: cachePath, size: fileSize}
	elem := fileCache.itemList.PushFront(cacheItem)
	fileCache.itemMap[filename] = elem
	fileCache.curSize += fileSize
}

// 读取缓存，如果命中则返回 true
func (fileCache *fileCache) ReadCache(w http.ResponseWriter, url *url.URL) bool {
	fileCache.lock.RLock()
	defer fileCache.lock.RUnlock()
	filename := strings.ReplaceAll(url.String(), "/", "@")
	// 转换路径并检查是否存在于缓存中
	elem, found := fileCache.itemMap[filename]
	if !found {
		return false // 缓存未命中
	}
	if log_cache {
		log.Println("命中缓存", filename)
	}
	// 移动缓存项到列表前端
	fileCache.itemList.MoveToFront(elem)

	cacheItem := elem.Value.(*cacheItem)
	data, err := os.ReadFile(cacheItem.filepath)
	if err != nil {
		return false
	}

	// 写入数据到响应
	w.Write(data)
	return true
}

// 移除最旧的缓存文件
func (fileCache *fileCache) removeOldest() {
	oldest := fileCache.itemList.Back()
	if oldest == nil {
		return
	}

	cacheItem := oldest.Value.(*cacheItem)
	fileCache.curSize -= cacheItem.size
	os.Remove(cacheItem.filepath) // 删除磁盘上的缓存文件
	delete(fileCache.itemMap, cacheItem.filename)
	fileCache.itemList.Remove(oldest)
	if log_cache {
		log.Println("删除缓存", cacheItem.filename)
	}
}
