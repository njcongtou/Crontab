package main

import (
	"Crontab/src/worker"
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
	// worker -config ./master.json
	// worker -h
	flag.StringVar(&confFile, "config", "./worker.json", "指定worker.json")
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
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}
	// 任务管理器
	if err = worker.InitJobMgr(); err != nil {
		return
	}

	for {
		time.Sleep(1 * time.Second)
	}

ERR:
	fmt.Println(err)
}
