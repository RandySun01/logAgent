package main
import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)
/*
@author RandySun
@create 2021-09-11-22:18
*/
func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil{
		fmt.Printf("connect to etcd failed err: %#v", err)
	}
	fmt.Println("connect to etcd success")
	defer cli.Close()
	// 添加数据
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	str := `[{"path":"G:/goproject/go/logAgent/demo/web.log","topic":"web_log"},{"path":"G:/goproject/go/logAgent/demo/randy.log","topic":"randy_log"},{"path":"G:/goproject/go/logAgent/demo/randy3.log","topic":"randy_log"}]`
	//str := `[{"path":"G:/goproject/go/logAgent/demo/web.log","topic":"web_log"},{"path":"G:/goproject/go/logAgent/demo/randy.log","topic":"randy_log"}]`

	_, err = cli.Put(ctx, "collect_log_conf", str)
	if err != nil{
		fmt.Printf("put to etcd failed, err: %#v", err)
	}
	cancel()


	// get 取数据
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, "collect_log_conf")
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
	}
	fmt.Println(resp.Kvs)
	//// del 取数据
	//ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	//delResponse, err := cli.Delete(ctx, "name")
	//cancel()
	//if err != nil {
	//	fmt.Printf("del from etcd failed, err:%v\n", err)
	//	return
	//}
	//fmt.Println("delete count: ",delResponse.Deleted)
}