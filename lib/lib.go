package lib

import (
	"github.com/cxpgo/golib/model"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)
var gConfig model.Config
var gVP     *viper.Viper
var gRedis  *Redis

type Map = map[string]interface{}
var Log *golibLog

var GORMMapPool map[string]*gorm.DB

var GORMDefaultPool *gorm.DB

func Init(path ...string)  {
	InitConfig(path...)
	InitLog(gConfig.Log)
	InitDBPool(gConfig.MySqlConfList)
	InitGormPool(gConfig.MySqlConfList)
	InitRedis(gConfig.RedisConfList["default"])
	//gRedis.Do("set","golib","ok")
	//gRedis.Do("SET", "k", "v")

}

func Destroy()  {
	Log.Close()
	CloseDB()
}
