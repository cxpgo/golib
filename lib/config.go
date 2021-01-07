package lib

import (
	"fmt"
	"github.com/cxpgo/golib/global"
	"github.com/spf13/viper"
	"os"
)

func InitConfig(path ...string) *viper.Viper{
	var config string
	if len(path) == 0 {
		//flag.StringVar(&config, "c", "", "choose config file.")
		//flag.Parse()
		if config == "" { // 优先级: 命令行 > 环境变量 > 默认值
			if configEnv := os.Getenv(global.ConfigEnv); configEnv == "" {
				config = global.ConfigFile
				fmt.Printf("您正在使用config的默认值,config的路径为%v\n", global.ConfigFile)
			} else {
				config = configEnv
				fmt.Printf("您正在使用GVA_CONFIG环境变量,config的路径为%v\n", config)
			}
		} else {
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%v\n", config)
		}
	} else {
		config = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%v\n", config)
	}

	v := viper.New()
	v.SetConfigFile(config)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	v.WatchConfig()

	if err := v.Unmarshal(&gConfig); err != nil {
		fmt.Println(err)
	}

	//fmt.Printf("Config : %+v \n",gConfig)

	return v
}
