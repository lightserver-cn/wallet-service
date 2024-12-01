package boot

import (
	"errors"
	"log"
	"os"
	"server/config"

	"gopkg.in/yaml.v3"
)

var (
	envCnfPath      = "/usr/local/config/config.yaml"
	envCnfLocalPath = "config/config.local.yaml"
)

func initConfig() error {
	env := os.Getenv("ENV")
	if env == "" {
		envCnfPath = envCnfLocalPath
	}

	err := unmarshalConfig(&config.Config, envCnfPath)
	if err != nil {
		return err
	}

	log.Printf("------ ENV:%s envCnfPath:%s Config:%v\n", env, envCnfPath, config.Config)

	return nil
}

func unmarshalConfig(conf any, filePath string) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return errors.New(err.Error() + filePath)
	}

	err = yaml.Unmarshal(bytes, conf)
	if err != nil {
		return errors.New(err.Error() + filePath)
	}

	return nil
}
