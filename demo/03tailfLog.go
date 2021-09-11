package main

/*
@author RandySun
@create 2021-09-11-14:46
*/

import (
	"fmt"
	"github.com/hpcloud/tail"
	"time"
)
// tailfile demo
func main() {
	fileName := `G:\goproject\go\logAgent\demo\randy.log`
	config := tail.Config{
		ReOpen:    true, // 打开文件
		Follow:    true, // 文件切割自动重新打开
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // Location读取文件的位置, Whence更加系统选择参数
		MustExist: false, // 允许日志文件不存在
		Poll:      true, // 轮询
	}
	// 打开文件读取日志
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		fmt.Println("tailfile file failed, err:", err)
		return
	}
	// 开始读取数据
	var (
		msg *tail.Line
		ok  bool
	)
	for {
		msg, ok = <-tails.Lines
		if !ok {
			fmt.Printf("tailfile file close reopen, filename:%s\n", tails.Filename)
			time.Sleep(time.Second) // 读取出错停止一秒
			continue
		}
		fmt.Println("msg:", msg.Text)
	}
}
