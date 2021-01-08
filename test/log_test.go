package test

import (
	"github.com/cxpgo/golib/lib"
	"os"
	"sync"
	"testing"
	"time"
)

var initOnce  sync.Once = sync.Once{}
func testInitOnce(logconf ...interface{}){
	initOnce.Do(func() {
		os.Chdir("..")
		lib.InitGolibConfig()
		lib.InitLog(lib.GConfig.Log)
	})
}

//测试日志打点
func TestDefaultLog(t *testing.T) {
	//os.Chdir("..")
	//ok,_ := utils.PathExists("logs")
	//fmt.Printf("logs is exist=%v",ok)
	testInitOnce()
	lib.Log.Error("TestDefaultLog")
	lib.Log.Infow(lib.Map{"infow:": "infow"}, lib.NewTrace())
	lib.Log.Infof("infof %s","is test")
	//Log.Close()
}

func TTestRotateLog(t *testing.T) {
	testInitOnce()
	//	GConfig.Log.Zap.Stdout = false
	//	GConfig.Log.Rotate.MaxSize = l1
	//	GConfig.Log.Rotate.MaxBackups = 3
	for{
		time.Sleep(time.Millisecond * 2)
		lib.Log.Error("TestRotateLog")
	}
	//Log.Close()

}
