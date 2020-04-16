package master

import (
	"net"
	"net/http"
	"time"
)

type ApiServer struct {
	httpServer *http.Server
}

var (
	// 单例对象
	// 需要被其他调用，就public
	G_apiServer *ApiServer
)

// 保存任务接口
func handleJobSave(w http.ResponseWriter, r *http.Request) {

}

// 初始化服务
// 需要被其他调用，就public
func InitApiServer() (err error) { // err 在这里定义住了，下面的return就不用再去返回err
	var (
		mux        *http.ServeMux
		listener   net.Listener
		httpServer *http.Server
	)

	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)

	if listener, err = net.Listen("tcp", ":8070"); err != nil {
		return
	}

	// 创建一个HTTP服务
	httpServer = &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      mux,
	}

	// 赋值单例
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}

	// 让这个http server开始工作
	// 启动了服务端
	go httpServer.Serve(listener)
	return
}
