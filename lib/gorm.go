package lib

import (
	"github.com/cxpgo/golib/model"
	"github.com/cxpgo/golib/utils/glog"

	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func InitGormPool(dbConfList map[string]*model.MySQLConf) error {
	//fmt.Printf("gorm %+v",dbConfList)
	GORMMapPool = map[string]*gorm.DB{}
	for confName, DbConf := range dbConfList {
		dataSourceName := getDataSourceNameByConfig(DbConf)

		newLogger := glog.GormLogNew(
			Log, // io writer
			glog.Config{
				SlowThreshold: time.Second, // Slow SQL threshold
				LogLevel:      logger.Info, // Log level
				Colorful:      false,       // Disable color
			},
		)

		dbgorm, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{Logger: newLogger, DisableForeignKeyConstraintWhenMigrating: true,})
		if err != nil {
			return err
		}

		sqlDB, err := dbgorm.DB()
		sqlDB.SetMaxIdleConns(DbConf.MaxIdleConn)
		sqlDB.SetMaxOpenConns(DbConf.MaxOpenConn)
		sqlDB.SetConnMaxLifetime(time.Second * time.Duration(DbConf.MaxConnLifeTime))

		//dbgorm.SingularTable(true)
		err = sqlDB.Ping()
		if err != nil {
			Log.Info("gorm connect is error")
			return err
		}

		GORMMapPool[confName] = dbgorm

	}

	if dbpool, err := GetGormPool("default"); err == nil {
		GORMDefaultPool = dbpool
	}

	Log.Info("===>Gorm Init Successful<===")

	return nil
}
func CloseGorm() {
	//for _, dbpool := range GORMMapPool {
	//
	//}
}

func GetGormPool(name string) (*gorm.DB, error) {
	if dbpool, ok := GORMMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get pool error")
}
