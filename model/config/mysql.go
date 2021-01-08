package config
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
