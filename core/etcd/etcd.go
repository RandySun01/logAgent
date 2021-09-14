package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"logAgent/common"
	"logAgent/core/tailfile"
	"time"
)

/*
@author RandySun
@create 2021-09-12-14:22
*/

// etcd相关操作
var (
	client *clientv3.Client
)

func Init(address []string) (err error) {
	client, err = clientv3.New(clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		fmt.Printf("connect to etcd failed, err:%v", err)
		return
	}
	return
}

// 拉取日志收集配置项的函数

func GetConf(key string) (collectEntryList []common.CollectEntry, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	resp, err := client.Get(ctx, key)
	if err != nil {
		logrus.Errorf("get conf from etcd by key:%s failed ,err:%v", key, err)
		return
	}
	if len(resp.Kvs) == 0 {
		logrus.Warningf("get len:0 conf from etcd by key:%s", key)
		return
	}
	ret := resp.Kvs[0]
	// ret.Value // json格式字符串
	fmt.Println("json格式字符串", string(ret.Value))
	err = json.Unmarshal(ret.Value, &collectEntryList)
	if err != nil {
		logrus.Errorf("json unmarshal failed, err:%v", err)
		return
	}
	return
}

// 监控etcd中日志收集项配置变化的函数

func WatchConf(key string) {
	for {
		// watch
		watchCh := client.Watch(context.Background(), key) // <-chan WatchResponse

		// 获取修改的指监控
		for wresp := range watchCh {
			for _, evt := range wresp.Events {
				// 获取被修改的key
				fmt.Printf("type:%s key:%s value: %s\n", evt.Type, evt.Kv.Key, evt.Kv.Value)
				var newConf []common.CollectEntry
				if evt.Type == clientv3.EventTypeDelete {
					// 如果是删除传空的文件配置
					logrus.Warningf("FBI warning: etcd delete the key!!!")
					tailfile.SendNewConf(   newConf)
				}
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil {
					logrus.Errorf("json unmarshal new conf failed, err: %v", err)
					continue
				}
				// 告诉tailfiel这个模块应该启用新的配置
				tailfile.SendNewConf(newConf) // 没有人接受，就是阻塞
			}
		}
	}

}
