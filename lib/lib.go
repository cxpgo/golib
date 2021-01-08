package lib

import (
	"database/sql"
	"github.com/cxpgo/golib/model/config"
	"gorm.io/gorm"
)

//definition common type
type Map = map[string]interface{}

var (
	Log *golibLog
	GConfig config.Config
	GGorm   *gorm.DB
	GDB     *sql.DB
	GRedis  *Redis
)

//Init config log mysql gorm redis
func InitAll(path ...string) {
	InitGolibConfig(path...)
	InitLog(GConfig.Log)
	InitDBPool(GConfig.MySqlConfList)
	InitGormPool(GConfig.MySqlConfList)
	InitRedis(GConfig.RedisConfList)
}
//Init config log
func InitOnlyLog(path ...string) {
	InitGolibConfig(path...)
	InitLog(GConfig.Log)

}
//Init config log mysql gorm
func InitLogDB(path ...string) {
	InitGolibConfig(path...)
	InitLog(GConfig.Log)
	InitDBPool(GConfig.MySqlConfList)
	InitGormPool(GConfig.MySqlConfList)
}
//Init config log redis
func InitLogRedis(path ...string) {
	InitGolibConfig(path...)
	InitLog(GConfig.Log)
	InitRedis(GConfig.RedisConfList)
}

//Close all
func Destroy() {
	CloseLog()
	CloseDB()
	CloseRedis()
}
