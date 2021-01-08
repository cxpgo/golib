package config
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
