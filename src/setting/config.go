package setting

import (
	"../service"

	"github.com/gentwolf-shen/gohelper/config"
	"github.com/gentwolf-shen/gohelper/dict"
	"github.com/gentwolf-shen/gohelper/logger"
	"github.com/gin-gonic/gin"
)

func Init() (config.Config, *gin.Engine) {
	cfg, err := config.LoadDefault()
	if err != nil {
		panic("load default config error: " + err.Error())
	}

	if err := dict.LoadDefault(); err != nil {
		panic("load dict error: " + err.Error())
		return cfg, nil
	}

	service.Auth.Init()

	var engine *gin.Engine
	if cfg.Web.IsDebug {
		engine = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		engine = gin.New()
	}

	logger.LoadDefault()
	router(engine)

	return cfg, engine
}
