package tailfile

import (
	"github.com/Shopify/sarama"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"logAgent/common"
	"logAgent/core/kafka"
	"strings"
	"time"
)

/*
@author RandySun
@create 2021-09-11-21:24
*/
var (
	TailObj *tail.Tail
)

// tail相关

type tailTask struct { // 多个日志文件对象
	path string
	topic string
	tailObj *tail.Tail
}

// 根据topic和path构造一个tailtask对象
func newTailTask(path, topic string)  (tt *tailTask){
	// 构造函数
	tt = &tailTask{
		path: path,
		topic: topic,
	}
	return
}

// 使用tail包打开日志文件准备读取日志
func (t *tailTask) Init() (err error){
	// 初始划tail对象
	config := tail.Config{
		ReOpen:    true, // 打开文件
		Follow:    true, // 文件切割自动重新打开
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2}, // Location读取文件的位置, Whence更加系统选择参数
		MustExist: false, // 允许日志文件不存在
		Poll:      true, // 轮询
	}

	t.tailObj, err = tail.TailFile(t.path, config)
	return
}
// 实际读取日志往kafka里面发送数据的方法
func (t *tailTask) run(){
	// logfile --> TailObj --> log --> Client --> kafka
	logrus.Infof("collect for path: %s is running....", t.path)
	for  {
		//循环读取数据
		line, ok := <- t.tailObj.Lines // chan tail.Line
		if !ok{
			logrus.Warn("tail file close reopen, path:%s\n", t.path)
			time.Sleep(time.Second) // 读取出错等一秒
			continue
		}

		// 如果是空行就略过
		//fmt.Printf("%#v\n", line.Text)
		if len(strings.Trim(line.Text,"\r")) == 0 {
			logrus.Info("出现空行拉,直接跳过...")
			continue
		}
		// 利用通道将同步的代码改为异步的
		// 把读出来的一行日志包装秤kafka里面的msg类型
		msg:=&sarama.ProducerMessage{}
		msg.Topic = t.topic // 每一个tailObj自己的topic
		msg.Value = sarama.StringEncoder(line.Text)
		// 将消息发送到通道中
		kafka.ToMsgChan(msg)

	}
}

// 初始化tailTask,为每一个日志文件构造一个单独的tailTask
func  Init(allEtcdConf []common.CollectEntry) (err error) {
	// allEtcdConf里面的内容: [{"path":"G:/goproject/go/logAgent/demo/web.log","topic":"web_log"},{"path":"G:/goproject/go/logAgent/demo/randy.log","topic":"randy_log"}]
	// 针对每一个日志收集项创建一个taskObj

	// 打开文件读取日志
	for _, conf := range allEtcdConf{

		tt := newTailTask(conf.Path, conf.Topic) // 创建一个日志收集任务
		err := tt.Init() // 去打开日志文件准备去读
		if err != nil{
			logrus.Errorf("create tailObj for path:%s, failed, err: %v", conf.Path, err)
		}
		logrus.Infof("create atail task for path: %s success", conf.Path)
		// 启一个后台的goroutine去收集日志
		go tt.run()
	}

	return
}