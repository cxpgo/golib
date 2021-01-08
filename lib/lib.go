package lib

import (
	"github.com/cxpgo/golib/model/config"
)

//definition common type
type Map = map[string]interface{}

var (
	GConfig config.Config
	GRedis  *Redis
)

//Init config log mysql gorm redis
func Init(path ...string) {
	InitGolibConfig(path...)
	InitLog(GConfig.Log)
	InitDBPool(GConfig.MySqlConfList)
	InitGormPool(GConfig.MySqlConfList)
	InitRedis(GConfig.RedisConfList)
}


//Close all
func Destroy() {
	CloseLog()
	CloseDB()
	CloseRedis()
}
