package lib
//
//import (
//	"bytes"
//	"io/ioutil"
//	"os"
//	"strings"
//	"time"
//
//	//dlog "github.com/cxpgo/golib/log"
//	//"github.com/iiswho666/gorm"
//	"github.com/spf13/viper"
//)
//
//type BaseConf struct {
//	DebugMode    string    `mapstructure:"debug_mode"`
//	TimeLocation string    `mapstructure:"time_location"`
//	Log          LogConfig `mapstructure:"log"`
//	Base         struct {
//		DebugMode    string `mapstructure:"debug_mode"`
//		TimeLocation string `mapstructure:"time_location"`
//	} `mapstructure:"base"`
//}
//
//type LogConfRotateFile struct {
//	RotateOn   bool   `mapstructure:"rotate_on"`
//	CommonPath string `mapstructure:"common_path"`
//	ErrorPath  string `mapstructure:"error_path"`
//	SqlPath    string `mapstructure:"sql_path"`
//	MaxSize    int    `mapstructure:"max_size"`
//	MaxBackups int    `mapstructure:"max_backups"`
//	MaxAge     int    `mapstructure:"max_age"`
//	Compress   bool   `mapstructure:"compress"`
//	LocalTime   bool   `mapstructure:"localTime"`
//}
//
//type LogConfZapEncoder struct {
//	ZapOn         string `mapstructure:"zap_on"`
//	Development   bool `mapstructure:"development"`
//	Stdout        bool `mapstructure:"stdout"`
//	Level         string `mapstructure:"level"`
//	TimeKey       string `mapstructure:"time_key"`
//	LevelKey      string `mapstructure:"level_key"`
//	NameKey       string `mapstructure:"name_key"`
//	CallerKey     string `mapstructure:"caller_key"`
//	MessageKey    string `mapstructure:"message_key"`
//	StacktraceKey string `mapstructure:"stack_trace_key"`
//}
//
//type LogConfig struct {
//	Level      string               `mapstructure:"log_level"`
//	RotateConf LogConfRotateFile    `mapstructure:"rotate_config"`
//	ZapConf    LogConfZapEncoder    `mapstructure:"zap_config"`
//}
//
////type MysqlMapConf struct {
////	List map[string]*MySQLConf `mapstructure:"list"`
////}
//
//type MySQLConf struct {
//	DriverName      string `mapstructure:"driver_name"`
//	DataSourceName  string `mapstructure:"data_source_name"`
//	MaxOpenConn     int    `mapstructure:"max_open_conn"`
//	MaxIdleConn     int    `mapstructure:"max_idle_conn"`
//	MaxConnLifeTime int    `mapstructure:"max_conn_life_time"`
//}
//
////type RedisMapConf struct {
////	List map[string]*RedisConf `mapstructure:"list"`
////}
//
////type RedisConf struct {
////	ProxyList    []string `mapstructure:"proxy_list"`
////	Password     string   `mapstructure:"password"`
////	Db           int      `mapstructure:"db"`
////	ConnTimeout  int      `mapstructure:"conn_timeout"`
////	ReadTimeout  int      `mapstructure:"read_timeout"`
////	WriteTimeout int      `mapstructure:"write_timeout"`
////}
//
////全局变量
//var ConfBase *BaseConf
//
////var ConfRedis *RedisConf
////var ConfRedisMap *RedisMapConf
//var ViperConfMap map[string]*viper.Viper
//
////获取基本配置信息
//func GetBaseConf() *BaseConf {
//	return ConfBase
//}
//
////InitBaseConf .
//func InitBaseConf(path string) error {
//	ConfBase = &BaseConf{}
//	err := ParseConfig(path, ConfBase)
//	if err != nil {
//		return err
//	}
//	//fmt.Printf("[INFO] ConfBase=%+v\n", ConfBase)
//	//fmt.Printf("[INFO] Log.Level=%v\n", ConfBase.Log.Level)
//	if ConfBase.DebugMode == "" {
//		if ConfBase.Base.DebugMode != "" {
//			ConfBase.DebugMode = ConfBase.Base.DebugMode
//		} else {
//			ConfBase.DebugMode = "debug"
//		}
//	}
//	if ConfBase.TimeLocation == "" {
//		if ConfBase.Base.TimeLocation != "" {
//			ConfBase.TimeLocation = ConfBase.Base.TimeLocation
//		} else {
//			ConfBase.TimeLocation = "Asia/Shanghai"
//		}
//	}
//	if ConfBase.Log.Level == "" {
//		ConfBase.Log.Level = "trace"
//	}
//	//配置 Zap 日志
//	//SetZapLogWithConf(ConfBase)
//
//	////配置日志
//	//logConf := dlog.LogConfig{
//	//	Level: ConfBase.Log.Level,
//	//	FW: dlog.ConfFileWriter{
//	//		On:              ConfBase.Log.FW.On,
//	//		LogPath:         ConfBase.Log.FW.LogPath,
//	//		RotateLogPath:   ConfBase.Log.FW.RotateLogPath,
//	//		WfLogPath:       ConfBase.Log.FW.WfLogPath,
//	//		RotateWfLogPath: ConfBase.Log.FW.RotateWfLogPath,
//	//	},
//	//	CW: dlog.ConfConsoleWriter{
//	//		On:    ConfBase.Log.CW.On,
//	//		Color: ConfBase.Log.CW.Color,
//	//	},
//	//}
//	//if err := dlog.SetupDefaultLogWithConf(logConf); err != nil {
//	//	panic(err)
//	//}
//	//dlog.SetLayout("2006-01-02T15:04:05.000")
//
//
//	return nil
//}
//
////
////func InitLogger(path string) error {
////	if err := dlog.SetupDefaultLogWithFile(path); err != nil {
////		panic(err)
////	}
////	dlog.SetLayout("2006-01-02T15:04:05.000")
////	return nil
////}
//
////func InitRedisConf(path string) error {
////	ConfRedis := &RedisMapConf{}
////	err := ParseConfig(path, ConfRedis)
////	if err != nil {
////		return err
////	}
////	ConfRedisMap = ConfRedis
////	return nil
////}
//
////初始化配置文件
//func InitViperConf() error {
//	f, err := os.Open(ConfEnvPath + "/")
//	if err != nil {
//		return err
//	}
//	//最多读取n个文件
//	fileList, err := f.Readdir(1024)
//	// log.Printf("[INFO] conf file list = %v\n", fileList)
//	if err != nil {
//		return err
//	}
//	for _, f0 := range fileList {
//		if !f0.IsDir() {
//			bts, err := ioutil.ReadFile(ConfEnvPath + "/" + f0.Name())
//			if err != nil {
//				return err
//			}
//			v := viper.New()
//			v.SetConfigType("toml")
//			v.ReadConfig(bytes.NewBuffer(bts))
//			pathArr := strings.Split(f0.Name(), ".")
//			if ViperConfMap == nil {
//				ViperConfMap = make(map[string]*viper.Viper)
//			}
//			ViperConfMap[pathArr[0]] = v
//		}
//	}
//	return nil
//}
//
////获取get配置信息
//func GetStringConf(key string) string {
//	keys := strings.Split(key, ".")
//	if len(keys) < 2 {
//		return ""
//	}
//	v, ok := ViperConfMap[keys[0]]
//	if !ok {
//		return ""
//	}
//	confString := v.GetString(strings.Join(keys[1:len(keys)], "."))
//	return confString
//}
//
////获取get配置信息
//func GetStringMapConf(key string) map[string]interface{} {
//	keys := strings.Split(key, ".")
//	if len(keys) < 2 {
//		return nil
//	}
//	v := ViperConfMap[keys[0]]
//	conf := v.GetStringMap(strings.Join(keys[1:len(keys)], "."))
//	return conf
//}
//
////获取get配置信息
//func GetConf(key string) interface{} {
//	keys := strings.Split(key, ".")
//	if len(keys) < 2 {
//		return nil
//	}
//	v := ViperConfMap[keys[0]]
//	conf := v.Get(strings.Join(keys[1:len(keys)], "."))
//	return conf
//}
//
////获取get配置信息
//func GetBoolConf(key string) bool {
//	keys := strings.Split(key, ".")
//	if len(keys) < 2 {
//		return false
//	}
//	v := ViperConfMap[keys[0]]
//	conf := v.GetBool(strings.Join(keys[1:len(keys)], "."))
//	return conf
//}
//
////获取get配置信息
//func GetFloat64Conf(key string) float64 {
//	keys := strings.Split(key, ".")
//	if len(keys) < 2 {
//		return 0
//	}
//	v := ViperConfMap[keys[0]]
//	conf := v.GetFloat64(strings.Join(keys[1:len(keys)], "."))
//	return conf
//}
//
////获取get配置信息
//func GetIntConf(key string) int {
//	keys := strings.Split(key, ".")
//	if len(keys) < 2 {
//		return 0
//	}
//	v := ViperConfMap[keys[0]]
//	conf := v.GetInt(strings.Join(keys[1:len(keys)], "."))
//	return conf
//}
//
////获取get配置信息
//func GetStringMapStringConf(key string) map[string]string {
//	keys := strings.Split(key, ".")
//	if len(keys) < 2 {
//		return nil
//	}
//	v := ViperConfMap[keys[0]]
//	conf := v.GetStringMapString(strings.Join(keys[1:len(keys)], "."))
//	return conf
//}
//
////获取get配置信息
//func GetStringSliceConf(key string) []string {
//	keys := strings.Split(key, ".")
//	if len(keys) < 2 {
//		return nil
//	}
//	v := ViperConfMap[keys[0]]
//	conf := v.GetStringSlice(strings.Join(keys[1:len(keys)], "."))
//	return conf
//}
//
////获取get配置信息
//func GetTimeConf(key string) time.Time {
//	keys := strings.Split(key, ".")
//	if len(keys) < 2 {
//		return time.Now()
//	}
//	v := ViperConfMap[keys[0]]
//	conf := v.GetTime(strings.Join(keys[1:len(keys)], "."))
//	return conf
//}
//
////获取时间阶段长度
//func GetDurationConf(key string) time.Duration {
//	keys := strings.Split(key, ".")
//	if len(keys) < 2 {
//		return 0
//	}
//	v := ViperConfMap[keys[0]]
//	conf := v.GetDuration(strings.Join(keys[1:len(keys)], "."))
//	return conf
//}
//
////是否设置了key
//func IsSetConf(key string) bool {
//	keys := strings.Split(key, ".")
//	if len(keys) < 2 {
//		return false
//	}
//	v := ViperConfMap[keys[0]]
//	conf := v.IsSet(strings.Join(keys[1:len(keys)], "."))
//	return conf
//}
