package model

type Config struct {
	System        System                `mapstructure:"system"`
	Log           Log                   `mapstructure:"log"`
	MySqlConfList map[string]*MySQLConf `mapstructure:"mysql"`
	RedisConfList map[string]*RedisConf `mapstructure:"redis"`
}

//系统配置
type System struct {
}

//ZAP 日志配置
type ZapConfig struct {
	On            bool   `mapstructure:"on"`
	MainPath      string `mapstructure:"main_path"`
	ErrorPath     string `mapstructure:"error_path"`
	Development   bool   `mapstructure:"development"`
	Stdout        bool   `mapstructure:"stdout"`
	Level         string `mapstructure:"level"`
	TimeKey       string `mapstructure:"time_key"`
	LevelKey      string `mapstructure:"level_key"`
	NameKey       string `mapstructure:"name_key"`
	CallerKey     string `mapstructure:"caller_key"`
	MessageKey    string `mapstructure:"message_key"`
	StacktraceKey string `mapstructure:"stack_trace_key"`
}

//文件转储配置
type RotateFileConf struct {
	On         bool   `mapstructure:"on"`
	FileName   string `mapstructure:"file_name"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
	LocalTime  bool   `mapstructure:"localTime"`
}

type Log struct {
	Zap    ZapConfig
	Rotate RotateFileConf
}

//Mysql 配置
type MySQLConf struct {
	Host            string `mapstructure:"host"`
	Port            string `mapstructure:"port"`
	Config          string `mapstructure:"config"`
	UserName        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DbName          string `mapstructure:"db_name"`
	MaxOpenConn     int    `mapstructure:"max_open_conn"`
	MaxIdleConn     int    `mapstructure:"max_idle_conn"`
	MaxConnLifeTime int    `mapstructure:"max_conn_life_time"`
}

//Redis 配置
type RedisConf struct {
	Host            string `mapstructure:"host"`
	Port            int `mapstructure:"port"`
	Password        string `mapstructure:"password"` // Password for AUTH.
	Db              int    `mapstructure:"db"`
	ReadTimeout     int    `mapstructure:"read_timeout"`
	WriteTimeout    int    `mapstructure:"write_timeout"`
	MaxIdle         int    `mapstructure:"conn_timeout"`     // Maximum number of connections allowed to be idle (default is 10)
	MaxActive       int    `mapstructure:"conn_timeout"`     // Maximum number of connections limit (default is 0 means no limit).
	ConnectTimeout  int    `mapstructure:"conn_timeout"`     // Dial connection timeout.
	IdleTimeout     int    `mapstructure:"idle_timeout"`     // Maximum idle time for connection (default is 10 seconds, not allowed to be set to 0)
	MaxConnLifetime int    `mapstructure:"max_conn_timeout"` // Maximum lifetime of the connection (default is 30 seconds, not allowed to be set to 0)
	TLS             bool   `mapstructure:"max_conn_timeout"` // Specifies the config to use when a TLS connection is dialed.
	TLSSkipVerify   bool   `mapstructure:"max_conn_timeout"` // Disables server name verification when connecting over TLS
}

//Host            string
//Port            int
//Db              int
//Pass            string        // Password for AUTH.
//MaxIdle         int           // Maximum number of connections allowed to be idle (default is 10)
//MaxActive       int           // Maximum number of connections limit (default is 0 means no limit).
//IdleTimeout     time.Duration // Maximum idle time for connection (default is 10 seconds, not allowed to be set to 0)
//MaxConnLifetime time.Duration // Maximum lifetime of the connection (default is 30 seconds, not allowed to be set to 0)
//ConnectTimeout  time.Duration // Dial connection timeout.
//TLS             bool          // Specifies the config to use when a TLS connection is dialed.
//TLSSkipVerify   bool          // Disables server name verification when connecting over TLS
