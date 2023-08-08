package global

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func path(path string, paths ...string) (string, error) {
	path, err := filepath.Abs(fmt.Sprintf("./%s", path))
	if err != nil {
		return "", err
	}

	_, existsErr := os.Stat(path)
	if os.IsNotExist(existsErr) {
		if err := os.MkdirAll(path, 0775); err != nil {
			return "", err
		}
	} else if existsErr != nil {
		return "", existsErr
	}

	if len(paths) > 0 {
		folder := strings.Join(paths, "/")
		path = fmt.Sprintf("%s/%s", path, folder)
	}

	return path, nil
}

func StoragePath(paths ...string) (string, error) {
	return path("storage", paths...)
}

func PublicPath(paths ...string) (string, error) {
	return path("public", paths...)
}

func ShouldStoragePath(path ...string) string {
	storage, err := StoragePath(path...)
	if err != nil {
		return ""
	}
	return storage
}

func ShouldPublicPath(path ...string) string {
	storage, err := PublicPath(path...)
	if err != nil {
		return ""
	}
	return storage
}
