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
}
