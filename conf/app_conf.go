package conf

import (
	"sync"
)

type Config struct {
	Env        string //开发环境
	KeyTimeOut int64  //验证码超时时间
	YPAppKey   string //云片appkey
	Dsn        string //mysql服务地址
}

var (
	Cfg     Config
	mutex   sync.Mutex //互斥锁
	declare sync.Once  //只执行一次
)

func Set(cfg Config) {
	mutex.Lock()
	Cfg.Env = setDefault(cfg.Env, "", "dev")
	Cfg.YPAppKey = setDefault(cfg.YPAppKey, "", "")
	Cfg.Dsn = setDefault(cfg.Dsn, "", "root:mysql@tcp(127.0.0.1:3306)/wdgl?charset=utf8mb4&parseTime=true&loc=Local")
	Cfg.KeyTimeOut = 30 * 60
	mutex.Unlock()
}

func setDefault(value, def, defValue string) string {
	switch value == def {
	case true:
		return defValue
	default:
		return value
	}
}
