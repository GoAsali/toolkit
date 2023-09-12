package config

import (
	"github.com/caarlos0/env/v8"
	"github.com/goasali/toolkit/global"
)

type IOPermission struct {
	Public  uint
	Private uint
}

type DiskPermission struct {
	File IOPermission
	Dir  IOPermission
}

type Disk struct {
	Driver        string
	Root          string
	Route         string
	Public        bool
	Permission    DiskPermission
	TemporaryLink bool
}

var fileSystem *FileSystem

type FileSystem struct {
	Default string `env:"FILESYSTEM_DRIVER"`
}

func GetFileSystem() (*FileSystem, error) {
	if fileSystem == nil {
		fileSystem = &FileSystem{}
		if err := env.Parse(fileSystem); err != nil {
			return nil, err
		}
		if fileSystem.Default == "" {
			fileSystem.Default = "local"
		}
	}
	return fileSystem, nil
}

func Disks() map[string]Disk {
	return map[string]Disk{
		"local": {
			Driver: "local",
			Root:   global.ShouldStoragePath(),
		},
		"upload": {
			Driver:        "local",
			Root:          global.ShouldStoragePath("uploads"),
			TemporaryLink: true,
		},
	}
}

func PublicDisks() map[string]Disk {
	disks := Disks()
	publicDisks := make(map[string]Disk)
	for key, disk := range disks {
		if disk.Public {
			publicDisks[key] = disk
		}
	}
	return publicDisks
}
