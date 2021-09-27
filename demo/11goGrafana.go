package main

import (
	"fmt"
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
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

var (
	lastNetIOStatTimeStamp int64    // 上一次获取IO数据的时间点
	lastNetInfo            *NetInfo //上一次的网络IO数据
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
	PartitionUsageStat map[string]*disk.UsageStat
}

// 存放网卡信息

type IOStat struct { //存放网卡速率
	BytesSent       uint64  `json:"bytes_sent"`   // number of bytes sent
	BytesRecv       uint64  `json:"bytes_recv"`   // number of bytes received
	PacketsSent     uint64  `json:"packets_sent"` // number of packets sent
	PacketsRecv     uint64  `json:"packets_recv"` // number of packets received
	BytesSentRate   float64 `json:"bytes_sent_rate"`
	BytesRecvRate   float64 `json:"bytes_recv_rate"`
	PacketsSentRate float64 `json:"packets_sent_rate"`
	PacketsRecvRate float64 `json:"packets_recv_rate"`
}
type NetInfo struct {
	//NetIOCountersStat map[string]*net.IOCountersStat
	NetIOCountersStat map[string]*IOStat
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

// 将网卡信息写入influxDB
func writesNetPoints(cli client.Client, sysInfo SysInfo) {
	netInfo := sysInfo.Data.(*NetInfo) // 类型转换

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	// 根据传入数据的类型，插入不同类型的数据
	for k, v := range netInfo.NetIOCountersStat {
		// 存放每一个节点
		tags := map[string]string{"name": k} // 每一个网卡存在tag中
		fields := map[string]interface{}{
			"bytes_sent_rate":   v.BytesSentRate,
			"bytes_recv_rate":   v.BytesRecvRate,
			"packets_sent_rate": v.PacketsSentRate,
			"packets_recv_rate": v.PacketsRecvRate,
		}
		// 添加数据
		pt, err := client.NewPoint("net", tags, fields, time.Now())
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
	log.Println("insert net success")
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

func getNetInfo1() {
	var sysInfo SysInfo

	var netInfo = &NetInfo{
		NetIOCountersStat: make(map[string]*IOStat, 18),
	}
	info, _ := net.IOCounters(true)
	currentTimeStamp := time.Now().Unix() // 获取当前网卡的时间
	for _, netIO := range info {
		//记录当前网卡字节数
		var ioStat = new(IOStat)
		ioStat.BytesSent = netIO.BytesSent     // 发送字节
		ioStat.BytesRecv = netIO.BytesRecv     // 接收字节
		ioStat.PacketsSent = netIO.PacketsSent // 发送包字节
		ioStat.PacketsRecv = netIO.PacketsRecv // 接收包字节

		// 记录
		netInfo.NetIOCountersStat[netIO.Name] = ioStat

		// 开始计算网卡相关速率
		if lastNetIOStatTimeStamp == 0 || lastNetInfo == nil {
			continue
		}
		// 计算时间间隔和速率
		interval := currentTimeStamp - lastNetIOStatTimeStamp
		// 计算速率
		bytesSentRate := float64(ioStat.BytesSent-lastNetInfo.NetIOCountersStat[netIO.Name].BytesSent) / float64(interval)
		bytesRecvRate := float64(ioStat.BytesRecv-lastNetInfo.NetIOCountersStat[netIO.Name].BytesRecv) / float64(interval)
		packetsSentRate := float64(ioStat.PacketsSent-lastNetInfo.NetIOCountersStat[netIO.Name].PacketsSent) / float64(interval)
		packetsRecvRate := float64(ioStat.PacketsRecv-lastNetInfo.NetIOCountersStat[netIO.Name].PacketsRecv) / float64(interval)

		// 记录速率
		ioStat.BytesSentRate = bytesSentRate
		ioStat.BytesRecvRate = bytesRecvRate
		ioStat.PacketsSentRate = packetsSentRate
		ioStat.PacketsRecvRate = packetsRecvRate

	}
	lastNetIOStatTimeStamp = currentTimeStamp // 更新上一次时间
	lastNetInfo = netInfo                     // 更新上一次发送包

	// 写入数据库中
	sysInfo.InfoType = NetInfoType
	sysInfo.Data = netInfo
	client := connInflux2()
	writesNetPoints(client, sysInfo)

}

func run(interval time.Duration) {
	tick := time.Tick(interval)
	for _ = range tick {
		getCpuInfo1()
		getMemInfo1()
		getDiskInfo1()
		getNetInfo1()
	}

}
func main() {
	// 每一秒写入数据
	run(time.Second)

}
