package boot

import (
	"server/config"
	"server/pkg/logger"
)

func initLog() error {
	logConf := config.Config.Log

	log, err := logger.InitZapLogger(
		logConf.FilePath,
		logConf.InfoFilename,
		logConf.WarnFilename,
		logConf.ErrFilename,
		logConf.FileExt,
		config.Config.AppName)
	if err != nil {
		return err
	}

	logger.Logger = log

	return nil
}
