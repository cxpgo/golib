package main

import (
	"fmt"
	"github.com/cxpgo/golib/lib"
	// "github.com/cxpgo/golib/stdlib"

)

func main() {
	fmt.Println("\n	欢迎使用 Golib	当前版本:V1.0.5\n ")

	lib.Init()
	defer lib.Destroy()


	//if err := lib.InitModule("./conf/dev/", []string{"base", "mysql", "redis"}); err != nil {
	//	log.Fatal(err)
	//}
	//defer lib.Destroy()
	////todo sth
	//lib.Log.Infow(lib.NewTraceByTag("Main"),lib.Map{"msg":"==>info", })
	//lib.Log.Error("==>error")
	//
	//for{
	//	time.Sleep(time.Millisecond * 2)
	//	lib.Log.Error("TestRotateLog")
	//}

	// stdlib.MainTest()

}
