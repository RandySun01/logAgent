package main

import (
	"fmt"
	"log"
	"net"
)

/*
@author RandySun
@create 2021-09-15-7:46
*/

// 获取ip
func main() {
	//ip, err := GetLocalIP()
	//fmt.Println(ip, err)
	ip:= GetOutboundIP()
	fmt.Println(ip)

}
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


// Get preferred outbound ip of this machine

func GetOutboundIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	fmt.Println(localAddr.String())
	return localAddr.IP.String()
}