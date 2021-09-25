package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"log"
	"time"
)

/*
@author RandySun
@create 2021-09-24-8:27
*/
const (
	CpuInfoType  = "cpu"
	MemInfoType  = "mem"
	DiskInfoType = "disk"
	NetInfoType  = "net"
)

// 存放系统信息结构体
type SysInfo struct {
	InfoType string
	Ip       string
	Data     interface{}
}

// 存放CPU使用
type CpuInfo struct {
	CpuPercent float64 `json:"cpu_percent"`
}

// 获取内存使用状态
type MemInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"used_percent"`
	Free        uint64  `json:"free"`
	Buffers     uint64  `json:"buffers"`
	Cached      uint64  `json:"cached"`
}

//  磁盘信息
// UsageStat 每一个磁盘分区的信息
type UsageStat struct {
	Path              string  `json:"path"`
	Fstype            string  `json:"fstype"`
	Total             uint64  `json:"total"`
	Free              uint64  `json:"free"`
	Used              uint64  `json:"used"`
	UsedPercent       float64 `json:"used_percent"`
	InodesTotal       uint64  `json:"inodes_total"`
	InodesUsed        uint64  `json:"inodes_used"`
	InodesFree        uint64  `json:"inodes_free"`
	InodesUsedPercent float64 `json:"inodes_used_percent"`
}

// PartitionStat 单个分区自己的信息
type PartitionStat struct {
	Device     string `json:"device"`
	Mountpoint string `json:"mountpoint"`
	Fstype     string `json:"fstype"`
	Opts       string `json:"opts"`
}

// DiskInfo 存放磁盘
type DiskInfo struct {
	PartitionUsageStat map[string]*disk.IOCountersStat
}

// 存放网卡信息

type NetInfo struct {
	NetIOCountersStat
}

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
func writesCpuPoints(cli client.Client, sysInfo SysInfo) {
	cpuInfo := sysInfo.Data.(*CpuInfo) // 类型转换

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	// 根据传入数据的类型，插入不同类型的数据
	tags := map[string]string{"cpu": "cpu0"}
	fields := map[string]interface{}{
		"cpu_percent": cpuInfo.CpuPercent,
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
	log.Println("insert cpu success")
}

// insert
func writesMemPoints(cli client.Client, sysInfo SysInfo) {
	memInfo := sysInfo.Data.(*MemInfo) // 类型转换

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	// 根据传入数据的类型，插入不同类型的数据
	tags := map[string]string{"mem": "mem"}
	fields := map[string]interface{}{
		"total":        int64(memInfo.Total),
		"available":    int64(memInfo.Available),
		"used":         int64(memInfo.Used),
		"used_percent": memInfo.UsedPercent,
		"free":         int64(memInfo.Free),
		"buffers":      int64(memInfo.Buffers),
		"cached":       int64(memInfo.Cached),
	}
	// 添加数据
	pt, err := client.NewPoint("memory", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// 写入数据
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert memory success")
}

// insert
func writesDiskPoints(cli client.Client, sysInfo SysInfo) {
	diskInfo := sysInfo.Data.(*DiskInfo) // 类型转换

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	// 根据传入数据的类型，插入不同类型的数据
	for k, v := range diskInfo.PartitionUsageStat {
		// 存放每一个节点
		tags := map[string]string{"path": k}
		fields := map[string]interface{}{
			"total":               int64(v.Total),
			"free":                int64(v.Free),
			"used":                int64(v.Used),
			"used_percent":        v.UsedPercent,
			"inodes_total":        int64(v.InodesTotal),
			"inodes_used":         int64(v.InodesUsed),
			"inodes_free":         int64(v.InodesFree),
			"inodes_used_percent": v.InodesUsedPercent,
		}
		// 添加数据
		pt, err := client.NewPoint("disk", tags, fields, time.Now())
		if err != nil {
			log.Fatal(err)
		}
		bp.AddPoint(pt)
	}

	// 写入数据
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert disk success")
}

func getCpuInfo1() {
	var sysInfo SysInfo
	var cpuInfo = new(CpuInfo)
	// CPU使用率
	percent, _ := cpu.Percent(time.Second, false)
	//fmt.Printf("cpu percent:%v\n", percent[0])
	// 将cpu采集信息写入到到influxdb中
	cpuInfo.CpuPercent = percent[0]
	sysInfo.InfoType = CpuInfoType
	sysInfo.Data = cpuInfo
	client := connInflux2()
	writesCpuPoints(client, sysInfo)

}

// mem info
func getMemInfo1() {
	var sysInfo SysInfo
	var memInfo = new(MemInfo)

	info, _ := mem.VirtualMemory()
	//fmt.Printf("mem info:%v\n", info)

	// 将内存信息放到结构体中
	memInfo.Total = info.Total
	memInfo.Available = info.Available
	memInfo.Used = info.Used
	memInfo.UsedPercent = info.UsedPercent
	memInfo.Free = info.Free
	memInfo.Buffers = info.Buffers
	memInfo.Cached = info.Cached
	// 存放内存
	sysInfo.InfoType = MemInfoType
	sysInfo.Data = memInfo
	client := connInflux2()

	writesMemPoints(client, sysInfo)

}

// disk info
func getDiskInfo1() {
	var sysInfo SysInfo

	var diskInfo = &DiskInfo{
		PartitionUsageStat: make(map[string]*disk.UsageStat, 18),
	}
	// 获取磁盘信息
	parts, err := disk.Partitions(true)
	if err != nil {
		fmt.Printf("get Partitions failed, err:%v\n", err)
		return
	}
	for _, part := range parts {
		//fmt.Printf("part:%v\n", part.String())
		// 获取每一个分区的信息
		diskInfoUsage, _ := disk.Usage(part.Mountpoint) // 传挂载点进去

		//fmt.Printf("disk info:used:%v free:%v\n", diskInfoUsage.UsedPercent, diskInfoUsage.Free)
		diskInfo.PartitionUsageStat[part.Mountpoint] = diskInfoUsage
	}

	sysInfo.InfoType = DiskInfoType
	sysInfo.Data = diskInfo
	client := connInflux2()
	writesDiskPoints(client, sysInfo)
}

func run(interval time.Duration) {
	tick := time.Tick(interval)
	for _ = range tick {
		getCpuInfo1()
		getMemInfo1()
		getDiskInfo1()
	}

}
func main() {
	// 每一秒写入数据
	run(time.Second)

}
