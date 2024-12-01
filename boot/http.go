package boot

import (
	"github.com/gin-gonic/gin"

	"log"

	"server/config"
	"server/pkg/dal"
	"server/pkg/logger"
	"server/router"
)

func initHTTP() error {
	gin.SetMode(config.Config.AppMode)

	engine := gin.Default()

	router.Router(engine, dal.CustomDal.DB, logger.Logger)

	log.Printf("start api server, address: %s, version: %s \n", config.Config.APIAddr, config.Config.AppVersion)

	err := engine.Run(config.Config.APIAddr)
	if err != nil {
		log.Fatalln("api run failed", config.Config.APIAddr, err.Error())
	}

	return nil
}
