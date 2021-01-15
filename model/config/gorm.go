package config

type Gorm struct {
	SqlLogLever   int `mapstructure:"sql_log_lever" json:"sqlLogLever"`
	TablePrefix string `mapstructure:"table_prefix" json:"tablePrefix" `
	SingularTable bool `mapstructure:"singular_table" json:"singularTable" `
}
