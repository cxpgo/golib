package config

type Config struct {
	System        System                `mapstructure:"system"`
	Http          Http                  `mapstructure:"http"`
	Log           Log                   `mapstructure:"log"`
	MySqlConfList map[string]*MySQLConf `mapstructure:"mysql"`
	RedisConfList map[string]*RedisConf `mapstructure:"redis"`
	Gorm          Gorm                  `mapstructure:"Gorm"`
}
