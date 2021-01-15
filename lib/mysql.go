package lib

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/cxpgo/golib/model/config"
	"go.uber.org/zap"
	"reflect"
	"regexp"
	"strconv"
	"time"
	"unicode"

	"gorm.io/gorm"
)

var DBMapPool map[string]*sql.DB

func InitDBPool(dbConfList map[string]*config.MySQLConf) error {
	DBMapPool = map[string]*sql.DB{}

	for confName, DbConf := range dbConfList {
		//
		dataSourceName := GetDataSourcePathByConfig(DbConf)
		dbpool, err := sql.Open("mysql", dataSourceName)
		if err != nil {
			Log.Errorw(map[string]interface{}{"err": err,}, NewTrace())
			return err
		}
		dbpool.SetMaxOpenConns(DbConf.MaxOpenConn)
		dbpool.SetMaxIdleConns(DbConf.MaxIdleConn)
		dbpool.SetConnMaxLifetime(time.Duration(DbConf.MaxConnLifeTime) * time.Second)
		err = dbpool.Ping()

		if err != nil {
			fmt.Println("mysql ping err====")
			Log.Infow(Map{"err": err,}, NewTrace())
			return err
		}

		DBMapPool[confName] = dbpool
	}

	if dbpool, err := GetGormPool("default"); err == nil {
		GGorm = dbpool
	}
	Log.Info("===>Mysql Init Successful<===")
	return nil
}
//default  return default config
func GetDataSourcePathByConfig(dbConf *config.MySQLConf) string {
	//"root:root3@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Local"
	//var build strings.Builder
	//build.WriteString(dbConf.UserName)
	//build.WriteString(":")
	//build.WriteString(dbConf.Password)
	//build.WriteString("@tcp(")
	//build.WriteString(dbConf.Host)
	//build.WriteString(")/")
	//build.WriteString(dbConf.DbName)
	//build.WriteString("?")
	//build.WriteString(dbConf.Config)
	//str := build.String()
	//str := dbConf.UserName + ":" + dbConf.Password + "@tcp(" + dbConf.Host + ")/" + dbConf.DbName + "?" + dbConf.Config
	if dbConf == nil{
		dbConf = GConfig.MySqlConfList["default"]
	}
	//fmt.Printf("dbConf=%+v\n",dbConf)
	str := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", dbConf.UserName, dbConf.Password, dbConf.Host,dbConf.Port, dbConf.DbName, dbConf.Config)
	return str
}

func GetDBPool(name string) (*sql.DB, error) {
	if dbpool, ok := DBMapPool[name]; ok {
		return dbpool, nil
	}
	return nil, errors.New("get pool error")
}

func CloseDB() error {
	for _, dbpool := range DBMapPool {
		dbpool.Close()
	}
	return nil
}

func DBPoolLogQuery(trace *TraceContext, sqlDb *sql.DB, query string, args ...interface{}) (*sql.Rows, error) {
	//startExecTime := time.Now()
	rows, err := sqlDb.Query(query, args...)
	//endExecTime := time.Now()
	if err != nil {
		//Log.TagError(trace, "_com_mysql_success", map[string]interface{}{
		//	"sql":       query,
		//	"bind":      args,
		//	"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExecTime).Seconds()),
		//})
	} else {
		//Log.TagInfo(trace, "_com_mysql_success", map[string]interface{}{
		//	"sql":       query,
		//	"bind":      args,
		//	"proc_time": fmt.Sprintf("%f", endExecTime.Sub(startExecTime).Seconds()),
		//})
	}
	return rows, err
}

//mysql日志打印类
// Logger default logger
type MysqlGormLogger struct {
	//gorm.Logger
	//Trace *TraceContext
}

// Print format & print log
func (logger *MysqlGormLogger) Print(values ...interface{}) {
	message := logger.LogFormatter(values...)
	if message["level"] == "sql" {
		//Log.TagInfo(logger.Trace, "_com_mysql_success", message)
	} else {
		//Log.TagInfo(logger.Trace, "_com_mysql_failure", message)
	}
}

// LogCtx(true) 时会执行改方法
func (logger *MysqlGormLogger) CtxPrint(s *gorm.DB, values ...interface{}) {
	//ctx, ok := s.GetCtx()
	//trace := NewTrace()
	//if ok {
	//	trace = ctx.(*TraceContext)
	//}
	message := logger.LogFormatter(values...)
	if message["level"] == "sql" {
		//Log.TagInfo(trace, "_com_mysql_success", message)
	} else {
		//Log.TagInfo(trace, "_com_mysql_failure", message)
	}
}

func (logger *MysqlGormLogger) LogFormatter(values ...interface{}) (messages map[string]interface{}) {
	if len(values) > 1 {
		var (
			sql             string
			formattedValues []string
			level           = values[0]
			currentTime     = logger.NowFunc().Format("2006-01-02 15:04:05")
			source          = fmt.Sprintf("%v", values[1])
		)

		messages = map[string]interface{}{"level": level, "source": source, "current_time": currentTime}

		if level == "sql" {
			// duration
			//messages = append(messages, fmt.Sprintf("%.2fms", float64(values[2].(time.Duration).Nanoseconds() / 1e4) / 100.0))
			messages["proc_time"] = fmt.Sprintf("%fs", values[2].(time.Duration).Seconds())
			// sql

			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format("2006-01-02 15:04:05")))
					} else if b, ok := value.([]byte); ok {
						if str := string(b); logger.isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formattedValues = append(formattedValues, "NULL")
				}
			}

			// differentiate between $n placeholders or else treat like ?
			if regexp.MustCompile(`\$\d+`).MatchString(values[3].(string)) {
				sql = values[3].(string)
				for index, value := range formattedValues {
					placeholder := fmt.Sprintf(`\$%d([^\d]|$)`, index+1)
					sql = regexp.MustCompile(placeholder).ReplaceAllString(sql, value+"$1")
				}
			} else {
				formattedValuesLength := len(formattedValues)
				for index, value := range regexp.MustCompile(`\?`).Split(values[3].(string), -1) {
					sql += value
					if index < formattedValuesLength {
						sql += formattedValues[index]
					}
				}
			}

			messages["sql"] = sql
			if len(values) > 5 {
				messages["affected_row"] = strconv.FormatInt(values[5].(int64), 10)
			}
		} else {
			messages["ext"] = values
		}
	}
	return
}

func (logger *MysqlGormLogger) NowFunc() time.Time {
	return time.Now()
}

func (logger *MysqlGormLogger) isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}

type GormLogger struct{}

// Print - Log Formatter
func (*GormLogger) Print(v ...interface{}) {
	switch v[0] {
	case "sql":
		zapLog.Info(
			"sql",
			zap.String("module", "gorm"),
			zap.String("type", "sql"),
			zap.Any("src", v[1]),
			zap.Any("duration", v[2]),
			zap.Any("sql", v[3]),
			zap.Any("values", v[4]),
			zap.Any("rows_returned", v[5]),
		)
	case "log":
		zapLog.Info("log", zap.Any("gorm", v[2]))
	}
}
