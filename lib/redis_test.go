package lib

import "testing"

func Test_Redis_Conn(t *testing.T) {
	testInitOnce()
	InitRedis(gConfig.RedisConfList["default"])
	_,err := gRedis.Do("Ping")
	if err != nil{
		t.Fatal(err)
	}
}
