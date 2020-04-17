package master

import (
	"Crontab/src/master/common"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/coreos/etcd/clientv3"
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
// POST job={"name": "job1", "command": "echo hello", "cronExpr": "* * * * *"}
func handleJobSave(w http.ResponseWriter, r *http.Request) {
	// 1. 解析post表单, 解析需要耗费cpu，http server不主动帮解析。需要手动去调用解析
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)

	if err = r.ParseForm(); err != nil {
		goto ERR
	}

	// 2. 取表单中的job字段
	postJob = r.PostForm.Get("job")

	// 3. 反序列化job
	if err = json.Unmarshal([]byte(postJob), &job); err != nil {
		goto ERR
	}

	// 4. 保存到etcd
	if oldJob, err = G_jobMgr.SaveJob(&job); err != nil {
		goto ERR
	}

	// 5. 返回正常应答 {{"errno": 0, "msg":"", "data": {...}})
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		w.Write(bytes)
	}

	return

ERR:
	// 6. 返回异常应答
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}
}

// 删除任务接口
// POST /job/delete name=job1
func handleJobDelete(w http.ResponseWriter, r *http.Request) {
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)

	// POST: a=1&b=2&c=3
	if err = r.ParseForm(); err != nil {
		goto ERR
	}

	// 删除的任务名
	name = r.PostForm.Get("name")

	// 去删除任务
	if oldJob, err = G_jobMgr.DeleteJob(name); err != nil {
		goto ERR
	}

	// 正常应答
	if bytes, err = common.BuildResponse(0, "success", oldJob); err == nil {
		w.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}
}

// 保存任务
func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {
	// 把任务保存到/cron/jobs/任务名 -> json
	var (
		jobKey    string
		jobValue  []byte
		putResp   *clientv3.PutResponse
		oldJobObj common.Job
	)

	jobKey = common.JOB_SAVE_DIR + job.Name
	if jobValue, err = json.Marshal(job); err != nil {
		return // error 会在返回值带回去
	}

	// 保存到etcd
	if putResp, err = jobMgr.kv.Put(context.TODO(), jobKey, string(jobValue), clientv3.WithPrevKV()); err != nil {
		return
	}

	// 如果是更新，那么返回旧值
	if putResp.PrevKv != nil {
		// 对旧值做一个反序列化
		if err = json.Unmarshal(putResp.PrevKv.Value, &oldJobObj); err != nil {
			err = nil
			return
		}
		oldJob = &oldJobObj
	}

	return
}

// 列举所有crontab任务
func handleJobList(w http.ResponseWriter, r *http.Request) {
	var (
		jobList []*common.Job
		bytes   []byte
		err     error
	)

	if jobList, err = G_jobMgr.ListJobs(); err != nil {
		goto ERR
	}

	// 正常应答
	if bytes, err = common.BuildResponse(0, "success", jobList); err == nil {
		w.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}

}

// 强制杀死某个任务
func handleJobKill(w http.ResponseWriter, r *http.Request) {
	var (
		err   error
		name  string
		bytes []byte
	)

	if err = r.ParseForm(); err != nil {
		goto ERR

	}

	name = r.PostForm.Get("name")

	// 杀死任务
	if err = G_jobMgr.KillJobs(name); err != nil {
		goto ERR
	}

	// 正常应答
	if bytes, err = common.BuildResponse(0, "success", nil); err == nil {
		w.Write(bytes)
	}
	return

ERR:
	if bytes, err = common.BuildResponse(-1, err.Error(), nil); err == nil {
		w.Write(bytes)
	}

}

// 初始化服务
// 需要被其他调用，就public
func InitApiServer() (err error) { // err 在这里定义住了，下面的return就不用再去返回err
	var (
		mux           *http.ServeMux
		listener      net.Listener
		httpServer    *http.Server
		staticDir     http.Dir     // 静态文件根目录
		staticHandler http.Handler // 静态文件的HTTP回调
	)

	// 配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

	// /index.html

	// 静态文件目录
	staticDir = http.Dir(G_config.Webroot)
	staticHandler = http.FileServer(staticDir)            // 用了内置的路由库
	mux.Handle("/", http.StripPrefix("/", staticHandler)) // /index.html -> index.html. 先去掉了/进行转化. ./webroot/index.html 得到了磁盘的路径

	if listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort)); err != nil {
		return
	}

	// 创建一个HTTP服务
	httpServer = &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Millisecond,
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
