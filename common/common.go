package common

import (
	"fmt"
	"net"
)

/*
@author RandySun
@create 2021-09-12-20:30
*/

// 要收集的日志配置结构体

type CollectEntry struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}

// 获取主机ip信息

func GetLocalIP() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet) // 类型断言

		if !ok {
			continue
		}

		if ipAddr.IP.IsLoopback() {
			continue
		}

		if !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		fmt.Println(ipAddr.IP)
		return ipAddr.IP.String(), nil
	}
	return
}