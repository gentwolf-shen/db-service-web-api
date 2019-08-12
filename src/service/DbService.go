package service

import (
	"strings"

	"../entity"

	"github.com/gentwolf-shen/gohelper/database"
)

var (
	Db = &DbService{}
)

type DbService struct {
}

func (this *DbService) Query(sqlMessage *entity.SqlMessage) ([]map[string]string, error) {
	return database.Driver(sqlMessage.Database).Query(sqlMessage.Sql, sqlMessage.Params...)
}

func (this *DbService) Update(sqlMessage *entity.SqlMessage) (map[string]int64, error) {
	n, err := database.Driver(sqlMessage.Database).Update(sqlMessage.Sql, sqlMessage.Params...)
	rs := map[string]int64{"affectedRows": n}
	return rs, err
}

func (this *DbService) Insert(sqlMessage *entity.SqlMessage) (map[string]int64, error) {
	n, err := database.Driver(sqlMessage.Database).Insert(sqlMessage.Sql, sqlMessage.Params...)
	rs := map[string]int64{"lastInsertId": n}
	return rs, err
}

// 批量查询, 同一个数据库, 不同的SQL
func (this *DbService) BatchQuery(sqlMessages []entity.SqlMessage) ([][]map[string]string, error) {
	rs := make([][]map[string]string, len(sqlMessages))

	for i, sqlMessage := range sqlMessages {
		if strings.ToUpper(sqlMessage.Sql[0:6]) != "SELECT" {
			continue
		}

		rows, err := this.Query(&sqlMessage)
		if err != nil {
			return nil, err
		}
		rs[i] = rows
	}

	return rs, nil
}

// 同一个数据库, 可以是不同的SQL, 仅处理insert、update、delete操作
func (this *DbService) Transaction(sqlMessages []entity.SqlMessage) (map[string]bool, error) {
	tx, err := database.Driver(sqlMessages[0].Database).GetConn().Begin()
	if err != nil {
		return nil, err
	}

	for _, sqlMessage := range sqlMessages {
		if strings.ToUpper(sqlMessage.Sql[0:6]) == "SELECT" {
			continue
		}

		if _, err = tx.Exec(sqlMessage.Sql, sqlMessage.Params...); err != nil {
			break
		}
	}

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return map[string]bool{"result": true}, err
}

// 同一个数据库, 同一个SQL, 不同参数, 仅处理insert, update, delete操作
func (this *DbService) TransactionV2(sqlMessage *entity.BatchSqlMessage) (map[string]bool, error) {
	tx, err := database.Driver(sqlMessage.Database).GetConn().Begin()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Prepare(sqlMessage.Sql)
	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	for _, param := range sqlMessage.Params {
		if _, err = stmt.Exec(param...); err != nil {
			break
		}
	}

	if err != nil {
		_ = tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	_ = stmt.Close()
	return map[string]bool{"result": true}, err
}
