package config

type Config struct {
	System        System         `mapstructure:"system"`
	Log           Log                   `mapstructure:"log"`
	MySqlConfList map[string]*MySQLConf `mapstructure:"mysql"`
	RedisConfList map[string]*RedisConf `mapstructure:"redis"`
}

