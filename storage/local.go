package storage

import (
	"github.com/gin-gonic/gin"
	"github.com/goasali/toolkit/config"
	"github.com/goasali/toolkit/global"
	services "github.com/goasali/toolkit/temporary"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type AbstractLocal struct {
	path         string
	route        string
	temporaryUrl bool
}

func NewAbstractLocal(disk config.Disk) *AbstractLocal {
	path := disk.Root
	if path[len(path)-1] != '/' {
		path = path + "/"
	}
	return &AbstractLocal{path, disk.Route, disk.TemporaryLink}
}

func (l AbstractLocal) GetOnlineLink(path string, optionFunctions ...OnlineLinkOptionsFunc) string {
	var route string
	if l.temporaryUrl {
		option := GetOnlineLinkOption(optionFunctions)

		key := path
		if option.uid != 0 {
			key += strconv.Itoa(int(option.uid))
		}

		key = global.GetMD5Hash(key)

		du, _ := time.ParseDuration("1h")
		temp, err := services.NewTemporary()
		if err != nil {
			panic(err)
		}
		err = temp.GenerateFileLink(path, services.WithExpiration(du), services.WithUser(option.uid))
		if err != nil {
			panic(err)
		}
		route = "temp?key=" + key
	} else {
		route = strings.ReplaceAll(path, l.ShouldPath(""), "")
	}
	app, err := config.GetApp()
	if err != nil {
		return "/" + route
	}
	return app.GetUrl(l.route + route)
}

func (l AbstractLocal) ServeOnRoute(prefix string, router *gin.Engine) {
	router.Static("/"+prefix, l.ShouldPath("/"))
}

func (l AbstractLocal) SaveUploadFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	fullPath, err := l.Path(dst)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(filepath.Dir(fullPath), 0750); err != nil {
		return err
	}

	out, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (l AbstractLocal) Exists(file string) (bool, error) {
	fullPath, err := l.Path(file)
	if err != nil {
		return false, err
	}
	_, errExists := os.Stat(fullPath)
	result := os.IsExist(errExists)
	return result, nil
}

func (l AbstractLocal) Delete(file string) error {
	fullPath, err := l.Path(file)
	if err != nil {
		return err
	}
	return os.Remove(fullPath)
}

func (l AbstractLocal) Read(file string) ([]byte, error) {
	fullPath, err := l.Path(file)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(fullPath)
}

func (l AbstractLocal) Path(file string) (path string, err error) {
	path, err = filepath.Abs(l.path)
	if err != nil {
		return "", err
	}
	return path + "/" + file, nil
}

func (l AbstractLocal) ShouldPath(file string) string {
	re, err := l.Path(file)
	if err != nil {
		return file
	}
	return re
}

func (l AbstractLocal) Write(filePath string, data []byte) error {
	return os.WriteFile(l.path+filePath, data, 0775)
}

func (l AbstractLocal) Copy(source string, target string) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	return l.Write(target, data)
}
