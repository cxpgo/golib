package config

type Log struct {
	Zap    ZapConfig
	Rotate RotateFileConf
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
