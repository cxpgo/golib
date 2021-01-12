package config
type Http struct {
	Post int `mapstructure:"post" json:"post" `
	GinModel string `mapstructure:"gin_model" json:"ginModel" gorm:"type:varchar(100);not null;column:gin_model" `
	ReadTimeout int `mapstructure:"read_timeout" json:"readTimeout" `
	WriteTimeout int `mapstructure:"write_timeout" json:"writeTimeout" `
	MaxHeaderBytes int `mapstructure:"max_header_bytes" json:"maxHeaderBytes" `
}
//实现TableName接口，以达到结构体和表对应，如果不实现该接口，并未设置全局表名禁用复数，gorm会自动扩展表名为articles（结构体+s）
//实现该接口后，设置的默认表前缀，对该表失效
func (*Http) TableName() string {
	return "http"
}
