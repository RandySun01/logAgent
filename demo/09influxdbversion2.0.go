package main

import (
	"context"
	"fmt"
	client "github.com/influxdata/influxdb-client-go/v2"
	"log"
	"time"
)

/*
@author RandySun
@create 2021-09-17-8:58
*/

import (
	"github.com/influxdata/influxdb-client-go/v2"
)

var bucket = "RandySun"
var org = "RandySun"

func connInflux() client.Client {
	// 连接数据库客户端
	token := "iLgyKP7N4-oTGKSj-vGVD8w9p-tHJQ-24BNouCfb4HEtHlSU-GeOCZ0cCWE3RauoSZDmVHJuB7Rg71Xd2b22sQ=="
	url := "http://127.0.0.1:8086"
	client := influxdb2.NewClient(url, token)
	return client
}

// insert
func writesPoints(client client.Client, org, bucket string) {

	tags := map[string]string{"cpu": "ih-cpu"}
	fields := map[string]interface{}{
		"idle":   201.1,
		"system": 43.3,
		"user":   86.6,
	}
	// 创建数据点
	pt := influxdb2.NewPoint("cpu_usage", tags, fields, time.Now())
	// 获取写入数据客户端
	writeAPI := client.WriteAPIBlocking(org, bucket)

	// 写入数据
	writeAPI.WritePoint(context.Background(), pt)
	log.Println("insert success")
}

//// query
func queryDB(client client.Client, org string) (err error) {

	queryAPI := client.QueryAPI(org)

	result, err := queryAPI.Query(context.Background(), `from(bucket:"RandySun")
    |> range(start: -1h) 
    |> filter(fn: (r) => r._measurement == "cpu_usage")`)
	if err == nil {
		for result.Next() {
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			fmt.Printf("value: %v\n", result.Record().Value())
		}
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}
	} else {
		panic(err)
	}
	return nil
}

func main() {
	client := connInflux()
	fmt.Println(client)
	writesPoints(client, org, bucket)
	queryDB(client, "RandySun")
}

//func main() {
//	conn := connInflux()
//	fmt.Println(conn)
//
//	// insert
//	writesPoints(conn)
//
//	// 获取10条数据并展示
//	qs := fmt.Sprintf("SELECT * FROM %s LIMIT %d", "cpu_usage", 10)
//	res, err := queryDB(conn, qs)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	for _, row := range res[0].Series[0].Values {
//		for j, value := range row {
//			log.Printf("j:%d value:%v\n", j, value)
//		}
//	}
//}
