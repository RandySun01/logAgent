package main

import (
	"context"
	"fmt"
	"github.com/influxdata/influxdb-client-go/v2"
	"log"
	"time"
)

/*
@author RandySun
@create 2021-09-17-8:58
*/

var bucket = "RandySun"
var org = "RandySun"

func connInflux() influxdb2.Client {
	// 连接数据库客户端
	token := "iLgyKP7N4-oTGKSj-vGVD8w9p-tHJQ-24BNouCfb4HEtHlSU-GeOCZ0cCWE3RauoSZDmVHJuB7Rg71Xd2b22sQ=="
	url := "http://127.0.0.1:8086"
	client := influxdb2.NewClient(url, token)
	return client
}

// insert
func writesPoints(client influxdb2.Client, org, bucket string) {

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
func queryDB(client influxdb2.Client, org string) (err error) {

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
