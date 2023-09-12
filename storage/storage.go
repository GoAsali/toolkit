package storage

import (
	"github.com/gin-gonic/gin"
	"github.com/goasali/toolkit/config"
	"mime/multipart"
)

type OnlineLinkOptionsFunc func(options *OnlineLinkOptions)

type OnlineLinkOptions struct {
	uid uint
}

func WithUserId(uid uint) OnlineLinkOptionsFunc {
	return func(options *OnlineLinkOptions) {
		options.uid = uid
	}
}

type IStorage interface {
	Path(string) (string, error)
	ShouldPath(string) string
	Write(string, []byte) error
	Copy(string, string) error
	Exists(string) (bool, error)
	Delete(string) error
	Read(string) ([]byte, error)
	SaveUploadFile(file *multipart.FileHeader, dst string) error
	GetOnlineLink(path string, optionFunctions ...OnlineLinkOptionsFunc) string
	ServeOnRoute(prefix string, router *gin.Engine)
}

func GetOnlineLinkOption(optionFunctions []OnlineLinkOptionsFunc) OnlineLinkOptions {
	option := OnlineLinkOptions{uid: 0}
	for _, optionFunc := range optionFunctions {
		optionFunc(&option)
	}
	return option
}

func Disk(diskName string) IStorage {
	return DiskFromConfig(config.Disks()[diskName])
}

func DiskFromConfig(disk config.Disk) IStorage {
	if disk.Driver == "local" {
		return NewAbstractLocal(disk)
	}
	panic("Unknown disk specified")
}

// Default get default upload path.
func Default() IStorage {
	driver, err := config.GetFileSystem()
	if err != nil {
		panic("Error in load filesystem config")
	}
	return Disk(driver.Default)
}
