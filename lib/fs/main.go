// Package fs gathers all of the function that helps manipulate the filesystem easily
package fs

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/aminnairi/gouache/lib/logger"
)

func Files(fileOrDirectoryPath string) func(yield func(filePath string) bool) {
	return func(yield func(filePath string) bool) {
		fileOrDirectoryStat, fileOrDirectoryStatError := os.Stat(fileOrDirectoryPath)

		if fileOrDirectoryStatError != nil {
			logger.Fatal("Error while reading path "+fileOrDirectoryPath+":", fileOrDirectoryStatError.Error())
			return
		}

		if !fileOrDirectoryStat.IsDir() {
			yield(fileOrDirectoryPath)
			return
		}

		walkError := filepath.WalkDir(fileOrDirectoryPath, func(path string, entry fs.DirEntry, err error) error {
			if err != nil {
				logger.Fatal("Error accessing path "+path+":", err.Error())
				return nil
			}

			if entry.Type().IsRegular() {
				if !yield(path) {
					return filepath.SkipAll
				}
			}
			return nil
		})

		if walkError != nil && walkError != filepath.SkipAll {
			logger.Fatal("Error while walking though the directory "+fileOrDirectoryPath+":", walkError.Error())
		}
	}
}

func FileExist(path string) bool {
	fileInfo, error := os.Stat(path)

	if error != nil {
		return false
	}

	return fileInfo.Mode().IsRegular()
}

func FolderExists(path string) bool {
	folderStat, statError := os.Stat(path)

	if statError != nil {
		return false
	}

	return folderStat.IsDir()
}
