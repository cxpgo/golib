package test

import (
	"github.com/cxpgo/golib/lib"
	"testing"
)

func Test_Redis_Conn(t *testing.T) {
	testInitOnce()
	lib.InitRedis(lib.GConfig.RedisConfList)
	_,err := lib.GRedis.Do("Ping")
	if err != nil{
		t.Fatal(err)
	}
}
