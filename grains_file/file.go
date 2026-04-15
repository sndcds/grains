package grains_file

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"
)

type DirInfo struct {
	FileCount    int    `json:"fileCount"`
	DirSize      string `json:"dirSize"`
	LastModified string `json:"lastModified"`
}

type MultiDirInfo struct {
	Name         string `json:"name"`
	FileCount    int    `json:"file_count"`
	DirSize      string `json:"dir_dize"`
	LastModified string `json:"last_modified"`
}

func FileStats(dir string) (int, int64, time.Time) {
	var fileCount int
	var totalSize int64
	var latestModTime time.Time

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			info, err := d.Info()
			if err != nil {
				return err
			}

			fileCount++
			totalSize += info.Size()

			// Track latest modification time
			if info.ModTime().After(latestModTime) {
				latestModTime = info.ModTime()
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
		return -1, -1, time.Time{}
	}

	return fileCount, totalSize, latestModTime
}

func MultiDirStats(dirs []string) []MultiDirInfo {
	results := make([]MultiDirInfo, 0, len(dirs))

	for _, dir := range dirs {
		fileCount, dirSize, lastMod := FileStats(dir)

		lastModified := ""
		if !lastMod.IsZero() {
			lastModified = lastMod.Format("2006-01-02 15:04:05")
		}

		results = append(results, MultiDirInfo{
			Name:         dir,
			FileCount:    fileCount,
			DirSize:      HumanSize(dirSize),
			LastModified: lastModified,
		})
	}

	return results
}

func HumanSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}
