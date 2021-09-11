package tailfile

import (
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
)

/*
@author RandySun
@create 2021-09-11-21:24
*/
var (
	TailObj *tail.Tail
)
// tail相关
func  Init(filename string) (err error) {
	config := tail.Config{
		ReOpen:    true, // 打开文件
		Follow:    true, // 文件切割自动重新打开
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // Location读取文件的位置, Whence更加系统选择参数
		MustExist: false, // 允许日志文件不存在
		Poll:      true, // 轮询
	}
	// 打开文件读取日志
	TailObj, err = tail.TailFile(filename, config)
	if err != nil {
		logrus.Error("tailfile: create tailObj for path:%s failed, err:%v\n", filename, err)
		return
	}
	return
}