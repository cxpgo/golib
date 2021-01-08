package lib

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"github.com/fsnotify/fsnotify"
)


var DefaultConfigEnv  = "GOLIB_CONFIG"
var DefaultConfigFile = "./config/dev/golibConfig.toml"

func InitGolibConfig(path ...string) *viper.Viper {
	//GConfig = &configPath.Config{}
	var configPath string
	if len(path) == 0 {
		flag.StringVar(&configPath, "c", "", "choose configPath file.")
		flag.Parse()
		if configPath == "" { // 优先级: 命令行 > 环境变量 > 默认值
			if configEnv := os.Getenv(DefaultConfigEnv); configEnv == "" {
				configPath = DefaultConfigFile
				fmt.Printf("您正在使用config的默认值,config的路径为%v\n", DefaultConfigEnv)
			} else {
				configPath = configEnv
				fmt.Printf("您正在使用GVA_CONFIG环境变量,config的路径为%v\n", configPath)
			}
		} else {
			fmt.Printf("您正在使用命令行的-c参数传递的值,config的路径为%v\n", configPath)
		}
	} else {
		configPath = path[0]
		fmt.Printf("您正在使用func Viper()传递的值,config的路径为%v\n", configPath)
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error configPath file: %s \n", err))
	}
	v.WatchConfig()

	//fmt.Printf("config===%d\n",len(config))

	if err := v.Unmarshal(&GConfig); err != nil {
		fmt.Println(err)
	}

	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(&GConfig); err != nil {
			fmt.Println(err)
		}
	})
	//fmt.Printf("Config : %+v \n", GConfig)

	return v
}

func InitConfig(configPath string, configModel interface{}) *viper.Viper {
	v := viper.New()
	v.SetConfigFile(configPath)
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error configPath file: %s \n", err))
	}
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := v.Unmarshal(configModel); err != nil {
			fmt.Println(err)
		}
	})
	//fmt.Printf("config===%d\n",len(config))
	if err := v.Unmarshal(configModel); err != nil {
		fmt.Println(err)
	}

	return v

}
