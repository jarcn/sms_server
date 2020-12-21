package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"log"
	"os"
	"sms_server/conf"
	"strconv"
)

var mysqlClt *xorm.Engine

var mysqlConf conf.MysqlConfig

func init() {
	if prodUrl := os.Getenv("mysql.datasource.url"); prodUrl != "" {
		mysqlConf.Dsn = prodUrl
	} else {
		mysqlConf.Dsn = conf.Db["db1"].Dsn
	}
	if showSql := os.Getenv("sql.show"); showSql != "" {
		mysqlConf.ShowSql, _ = strconv.ParseBool(showSql)
	} else {
		mysqlConf.ShowSql = conf.Db["db1"].ShowSql
	}
	mysqlConf.MaxIdle = conf.Db["db1"].MaxIdle
	mysqlConf.ShowExecTime = conf.Db["db1"].ShowExecTime
	mysqlConf.MaxOpen = conf.Db["db1"].MaxOpen
	mysqlConf.DriverName = conf.Db["db1"].DriverName
	if nil == mysqlClt {
		var err error
		mysqlClt, err = xorm.NewEngine(mysqlConf.DriverName, mysqlConf.Dsn)
		if err != nil {
			log.Fatal(err)
		}
		mysqlClt.SetMaxIdleConns(mysqlConf.MaxIdle) //空闲连接
		mysqlClt.SetMaxOpenConns(mysqlConf.MaxOpen) //最大连接数
		mysqlClt.ShowSQL(mysqlConf.ShowSql)
		mysqlClt.ShowExecTime(mysqlConf.ShowExecTime)
	}
}
