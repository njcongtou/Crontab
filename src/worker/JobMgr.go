package worker

import (
	"Crontab/src/master/common"
	"context"
	"time"

	"github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
)

type JobMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr // Singleton
)

// 监听任务变化
func (jobMgr *JobMgr) watchJobs() (err error) {
	var (
		getResp *clientv3.GetResponse
		kvpair  *mvccpb.KeyValue
		job     *common.Job
	)

	// 1. get 一下/cron/jobs 目录下的所有任务，并且获得当前集群的revision
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	// 当前哪些任务
	for _, kvpair = range getResp.Kvs {
		// 序列化json成job
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			// TODO: 把这个job同步给scheduler调度携程

		}
	}

	// 2. 从该revision向后监听变化事件

	return
}

func InitJobMgr() (err error) {
	var (
		config clientv3.Config
		client *clientv3.Client
		kv     clientv3.KV
		lease  clientv3.Lease
	)

	// 初始化配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond,
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}

	// 得到kv和lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	// 赋值单例
	G_jobMgr = &JobMgr{
		client: client,
		kv:     kv,
		lease:  lease,
	}

	return
}
