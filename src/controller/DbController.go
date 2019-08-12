package controller

import (
	"../service"

	"github.com/gin-gonic/gin"
)

type DbController struct {
	Controller
}

func (this *DbController) Query(c *gin.Context) {
	this.showResult(c, func(c *gin.Context) (interface{}, error) {
		if sqlMessage, err := this.getParam(c); err != nil {
			return nil, err
		} else {
			return service.Db.Query(sqlMessage)
		}
	})
}

func (this *DbController) Update(c *gin.Context) {
	this.showResult(c, func(c *gin.Context) (interface{}, error) {
		if sqlMessage, err := this.getParam(c); err != nil {
			return nil, err
		} else {
			return service.Db.Update(sqlMessage)
		}
	})
}

func (this *DbController) Insert(c *gin.Context) {
	this.showResult(c, func(c *gin.Context) (interface{}, error) {
		if sqlMessage, err := this.getParam(c); err != nil {
			return nil, err
		} else {
			return service.Db.Insert(sqlMessage)
		}
	})
}

// 批量查询, 同一个数据库, 不同的SQL
func (this *DbController) BatchQuery(c *gin.Context) {
	this.showResult(c, func(c *gin.Context) (interface{}, error) {
		if sqlMessages, err := this.getParams(c); err != nil {
			return nil, err
		} else {
			return service.Db.BatchQuery(sqlMessages)
		}
	})
}

// 同一个数据库, 可以是不同的SQL, 仅处理insert, update, delete操作
func (this *DbController) Transaction(c *gin.Context) {
	this.showResult(c, func(c *gin.Context) (interface{}, error) {
		if sqlMessages, err := this.getParams(c); err != nil {
			return nil, err
		} else {
			return service.Db.Transaction(sqlMessages)
		}
	})
}

// 同一个数据库, 同一个SQL, 不同的参数, 仅处理insert, update, delete操作
func (this *DbController) TransactionV2(c *gin.Context) {
	this.showResult(c, func(c *gin.Context) (interface{}, error) {
		if sqlMessage, err := this.getParamBatch(c); err != nil {
			return nil, err
		} else {
			return service.Db.TransactionV2(sqlMessage)
		}
	})
}
