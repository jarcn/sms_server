package conf

type MysqlConfig struct {
	DriverName   string
	Dsn          string
	ShowSql      bool
	ShowExecTime bool
	MaxIdle      int
	MaxOpen      int
}

var Db = map[string]MysqlConfig{
	"db1": {
		DriverName:   "mysql",
		Dsn:          "root:mysql@tcp(127.0.0.1:3306)/wdgl?charset=utf8mb4&parseTime=true&loc=Local",
		ShowSql:      true,
		ShowExecTime: true,
		MaxIdle:      10,
		MaxOpen:      200,
	},
}
