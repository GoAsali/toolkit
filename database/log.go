package database

import (
	"context"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"io"
)

type Logger struct {
	logger.Interface
}

func (gLog Logger) LogMode(level logger.LogLevel) logger.Interface {
	if level == logger.Silent {
		log.SetOutput(io.Discard)
	} else {
		levels := map[logger.LogLevel]string{
			logger.Error: "error",
			logger.Info:  "info",
			logger.Warn:  "warn",
		}
		parseLevel, _ := log.ParseLevel(levels[level])
		log.SetLevel(parseLevel)
	}
	return gLog
}

func (gLog Logger) Info(_ context.Context, msg string, data ...interface{}) {
	log.Println(msg, data)
}

func (gLog Logger) Warn(_ context.Context, msg string, data ...interface{}) {
	log.Println(msg, data)
}

func (gLog Logger) Error(_ context.Context, msg string, data ...interface{}) {
	log.Println(msg, data)
}

func (Logger) Print(v ...interface{}) {

}
