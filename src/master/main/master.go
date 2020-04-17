package main

import (
	"Crontab/src/master"
	"flag"
	"fmt"
	"runtime"
	"time"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	// master -config ./master.json
	// master -h
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
	flag.Parse()
}

func initEnv() {
	// 发挥go语言最大性能，就要找出当前CPU核心数
	// 就设置多少个协程 goroutines
	runtime.GOMAXPROCS(runtime.NumCPU())
	//fmt.Println(numCPU)
}

func main() {
	var (
		err error
	)

	// 初始化命令行参数
	initArgs()

	// 初始化线程
	initEnv()

	// load config
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 任务管理器
	if err = master.InitJobMgr(); err != nil {
		return
	}

	// 启动ApiHTTP服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	for {
		time.Sleep(1 * time.Second)
	}

ERR:
	fmt.Println(err)
}
