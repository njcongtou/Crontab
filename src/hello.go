package main

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/etcd/clientv3"
)

func main() {
	var (
		config clientv3.Config
		client *clientv3.Client
		err    error
		kv     clientv3.KV
		//putResp *clientv3.PutResponse
		getResp *clientv3.GetResponse
	)

	config = clientv3.Config{
		Endpoints:   []string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	}
	if client, err = clientv3.New(config); err != nil {
		// handle error!
	}
	defer client.Close()

	kv = clientv3.NewKV(client)

	/*
	   if putResp, err = kv.Put(context.TODO(), "/cron/jobs/job3", "bye3", clientv3.WithPrevKV()); err != nil {
	       fmt.Println(err)
	   } else {

	       fmt.Println("Revision: ", putResp.Header.Revision)
	       if putResp.PrevKv != nil {
	           fmt.Printf("Revision: %s", putResp.PrevKv.Value)
	       }

	   }
	*/
	if getResp, err = kv.Get(context.TODO(), "/cron/jobs/job3", clientv3.WithPrevKV()); err != nil {
		fmt.Println(err)
	} else {

		fmt.Println(getResp.Kvs)
	}

}
