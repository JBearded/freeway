package logger

import (
	"freeway/common"
	"freeway/config"
	"os"

	log "github.com/sirupsen/logrus"
)

// Init 初始化日志框架
func Init(profile common.Profile) error {
	filePath := config.Get().Logger.Path
	if fileInfo, statErr := os.Stat(filePath); statErr != nil && fileInfo == nil {
		_, createErr := os.Create(filePath)
		if createErr != nil {
			log.Fatal("create log file error", filePath, createErr)
			return createErr
		}
	}
	file, err := os.OpenFile(config.Get().Logger.Path, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal("open or create log file error", err)
		return err
	}
	level, _ := log.ParseLevel(config.Get().Logger.Level)

	if profile == common.DevelopPorfile || profile == common.NoneProfile {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(file)
	}
	log.SetLevel(level)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
	return nil
}
