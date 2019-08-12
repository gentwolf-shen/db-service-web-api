package controller

import (
	"../entity"
	"../service"

	"github.com/gentwolf-shen/gohelper/ginhelper"
	"github.com/gin-gonic/gin"
)

type Controller struct {
}

// 取单条SQL语言
func (this *Controller) getParam(c *gin.Context) (*entity.SqlMessage, error) {
	sqlMessage := &entity.SqlMessage{}
	if err := c.BindJSON(sqlMessage); err != nil {
		return nil, err
	}

	dbName, err := service.Auth.CheckSql(c.GetString("appKey"), sqlMessage.Sql)
	if err == nil {
		sqlMessage.Database = dbName
	}
	return sqlMessage, err
}

// 取多条SQL语言
func (this *Controller) getParams(c *gin.Context) ([]entity.SqlMessage, error) {
	var sqlMessages []entity.SqlMessage
	if err := c.BindJSON(&sqlMessages); err != nil {
		return sqlMessages, err
	}

	appKey := c.GetString("appKey")
	for i, v := range sqlMessages {
		dbName, err := service.Auth.CheckSql(appKey, v.Sql)
		if err != nil {
			return nil, err
		}

		sqlMessages[i].Database = dbName
	}

	return sqlMessages, nil
}

// 取批量操作SQL语言 (SQL相同，参数不同)
func (this *Controller) getParamBatch(c *gin.Context) (*entity.BatchSqlMessage, error) {
	sqlMessage := &entity.BatchSqlMessage{}
	if err := c.ShouldBindJSON(sqlMessage); err != nil {
		return nil, err
	}

	dbName, err := service.Auth.CheckSql(c.GetString("appKey"), sqlMessage.Sql)
	if err == nil {
		sqlMessage.Database = dbName
	}

	return sqlMessage, err
}

func (this *Controller) showResult(c *gin.Context, fun func(c *gin.Context) (interface{}, error)) {
	if result, err := fun(c); err == nil {
		ginhelper.ShowSuccess(c, result)
	} else {
		ginhelper.ShowErrorMsg(c, 4001001, err.Error())
	}
}
