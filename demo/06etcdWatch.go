package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

/*
@author RandySun
@create 2021-09-12-10:21
*/

// watch: 监控etcd中的Key的变化(创建，更改，修改)
func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})

	if err  != nil{
		fmt.Printf("connect fo failed err: %#v", err)
	}
	defer cli.Close()
	// watch
	watchCh := cli.Watch(context.Background(), "name") // <-chan WatchResponse

	// 获取修改的指监控
	for wresp := range watchCh{
		for _, env := range wresp.Events{
			// 获取被修改的key
			fmt.Printf("type:%s key:%s value: %s\n", env.Type, env.Kv.Key, env.Kv.Value)

		}
	}

}