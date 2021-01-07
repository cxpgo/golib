package lib

import (
	"github.com/cxpgo/golib/model"
	"github.com/spf13/viper"
)
//definition common type
type Map = map[string]interface{}

//definition global var
var gConfig model.Config
var gVP     *viper.Viper
var gRedis  *Redis

//Init config log mysql gorm redis
func Init(path ...string)  {
	InitConfig(path...)
	InitLog(gConfig.Log)
	InitDBPool(gConfig.MySqlConfList)
	InitGormPool(gConfig.MySqlConfList)
	InitRedis(gConfig.RedisConfList)
}

//Close all
func Destroy()  {
	CloseLog()
	CloseDB()
	CloseRedis()
}
