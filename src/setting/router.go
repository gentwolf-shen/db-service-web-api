package setting

import (
	"../controller"

	"github.com/gin-gonic/gin"
)

func router(engine *gin.Engine) {
	engine.Use(gin.Recovery())
	engine.Use(Auth())

	db := engine.Group("/")
	{
		ctl := controller.DbController{}
		db.POST("/query", ctl.Query)
		db.POST("/update", ctl.Update)
		db.POST("/delete", ctl.Update)
		db.POST("/insert", ctl.Insert)
		db.POST("/batch/query", ctl.BatchQuery)
		db.POST("/v1/transaction", ctl.Transaction)
		db.POST("/v2/transaction", ctl.TransactionV2)
	}
}
