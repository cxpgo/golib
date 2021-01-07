package lib

import (
	"os"
	"sync"
	"testing"
	"time"
)

var initOnce  sync.Once = sync.Once{}
func testInitOnce(logconf ...interface{}){
	initOnce.Do(func() {
		os.Chdir("..")
		InitConfig()
		InitLog(gConfig.Log)
	})
}

//测试日志打点
func TestDefaultLog(t *testing.T) {
	//os.Chdir("..")
	//ok,_ := utils.PathExists("logs")
	//fmt.Printf("logs is exist=%v",ok)
	testInitOnce()
	Log.Error("TestDefaultLog")
	Log.Infow(Map{"infow:":"infow"},NewTrace())
	Log.Infof("infof %s","is test")
	//Log.Close()
}

func TTestRotateLog(t *testing.T) {
	testInitOnce()
	//	gConfig.Log.Zap.Stdout = false
	//	gConfig.Log.Rotate.MaxSize = l1
	//	gConfig.Log.Rotate.MaxBackups = 3
	for{
		time.Sleep(time.Millisecond * 2)
		Log.Error("TestRotateLog")
	}
	//Log.Close()

}
