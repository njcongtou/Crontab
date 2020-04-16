package main

import (
	"Crontab/src/master"
	"fmt"
	"runtime"
)

func initEnv() {
	// 发挥go语言最大性能，就要找出当前CPU核心数
	// 就设置多少个协程 goroutines
	numCPU := runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println(numCPU)
}

func main() {
	var (
		err error
	)

	// 初始化线程
	initEnv()

	// 启动ApiHTTP服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

ERR:
	fmt.Println(err)
}
