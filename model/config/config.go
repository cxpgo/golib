package config

type Config struct {
	Log           Log                   `mapstructure:"log"`
	MySqlConfList map[string]*MySQLConf `mapstructure:"mysql"`
	RedisConfList map[string]*RedisConf `mapstructure:"redis"`
}

