package entity

type (
	SqlItem struct {
		Database string `json:"database"`
		Sql      string `json:"sql" binding:"required"`
	}

	SqlMessage struct {
		SqlItem

		Params []interface{} `json:"params"`
	}

	BatchSqlMessage struct {
		SqlItem

		Params [][]interface{} `json:"params"`
	}

	AuthConfig struct {
		Secret   string   `json:"secret"`
		Database string   `json:"database"`
		Actions  []string `json:"actions"`
	}
)
