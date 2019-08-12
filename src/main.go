package main

import (
	"./setting"

	_ "github.com/Go-SQL-Driver/MySQL"
	_ "github.com/gentwolf-shen/gohelper/daemon"
	"github.com/gentwolf-shen/gohelper/database"
	"github.com/gentwolf-shen/gohelper/endless"
)

func main() {
	cfg, engine := setting.Init()

	if err := database.LoadFromConfig(cfg.Db); err != nil {
		panic(err)
	} else {
		defer database.CloseAll()
	}

	if err := endless.ListenAndServe(cfg.Web.Port, engine); err != nil {
		panic(err)
	}
}
