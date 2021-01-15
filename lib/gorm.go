package lib

import (
	"github.com/cxpgo/golib/model/config"
	"github.com/cxpgo/golib/utils/glog"
	"gorm.io/gorm/schema"

	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

var GORMMapPool map[string]*gorm.DB

func InitGormPool(dbConfList map[string]*config.MySQLConf) error {
	//fmt.Printf("gorm %+v",dbConfList)
	GORMMapPool = map[string]*gorm.DB{}

	for confName, DbConf := range dbConfList {
		dataSourceName := GetDataSourcePathByConfig(DbConf)
		newLogger := glog.GormLogNew(
			Log, // io writer
			glog.Config{
				SlowThreshold: time.Second,                                  // Slow SQL threshold
				LogLevel:      logger.LogLevel(GConfig.Gorm.SqlLogLever), 	// Log level
				Colorful:      false,                                        // Disable color
			},
		)

		//mysqlConfig := mysql.Config{
		//	DSN:                       dataSourceName,   // DSN data source name
		//	DefaultStringSize:         191,   // string 类型字段的默认长度
		//	DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		//	DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		//	DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		//	SkipInitializeWithVersion: false, // 根据版本自动配置
		//}
		dbgorm, err := gorm.Open(mysql.Open(dataSourceName), &gorm.Config{
		//dbgorm, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
			Logger: newLogger,
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: GConfig.Gorm.TablePrefix,   // 表名前缀，`Article` 的表名应该是 `it_articles`
				SingularTable: GConfig.Gorm.SingularTable, // 使用单数表名，启用该选项，此时，`Article` 的表名应该是 `it_article`
			},
			DisableForeignKeyConstraintWhenMigrating: true,})
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
		GGorm = dbpool
	}

	Log.Info("===>GGorm Init Successful<===")

	return nil
}
func CloseGorm() {
	//for _, dbpool := range GORMMapPool {
	//
	//}
}

func GetGormPool(name ...string) (*gorm.DB, error) {
	var realName string

	if len(name) > 0 {
		realName = name[0]
	} else {
		realName = "default"
	}
	if dbpool, ok := GORMMapPool[realName]; ok {
		return dbpool, nil
	}

	return nil, errors.New("get pool error")
}

