package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/cpu"
	"log"
	"time"
)

/*
@author RandySun
@create 2021-09-24-8:27
*/

func connInflux2() client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "admin",
		Password: "",
	})
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

// insert
func writesPoints2(cli client.Client, perent float64) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	// 添加多条数据
	tags := map[string]string{"cpu": "cpu0"}
	fields := map[string]interface{}{
		"cpu_percent": perent,
	}
	// 添加数据
	pt, err := client.NewPoint("cpu_percent", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// 写入数据
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert success")
}
func getCpuInfo1() {

	// CPU使用率

	percent, _ := cpu.Percent(time.Second, false)
	fmt.Printf("cpu percent:%v\n", percent[0])
	// 将cpu采集信息写入到到influxdb中

	client := connInflux2()
	writesPoints2(client, percent[0])

}
func main() {
	// 每一秒写入数据
	for {
		getCpuInfo1()
		time.Sleep(time.Second)

	}

}
