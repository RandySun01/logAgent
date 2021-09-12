package main

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
	"logAgent/core/etcd"
	"logAgent/core/kafka"
	"logAgent/core/tailfile"
	"strings"
	time "time"
)

/*
@author RandySun
@create 2021-09-11-11:23
*/

// 日志收集客户端
// 类似的开源项目还有filebeat
// 收集指定目录下的日志文件,发送到kafka中

// 整个logagent的配置结构体
type Config struct {
	KafkaConfig `ini:"kafka"`
	CollectConfig `ini:"collect"`
	EtcdConfig `ini:"etcd"`
}
type KafkaConfig struct {
	Address string `ini:"address"`
	ChangeSize int64 `ini:"chan_size"`
}
type CollectConfig struct {
	LogFilePath string `ini:"logfile_path"`
}

type EtcdConfig struct {
	Address string `ini:"address"`
	CollectKey string `ini:"collect_key"`
}

// run 真正的业务逻辑
func  run()(err error){
	// logfile --> TailObj --> log --> Client --> kafka
	for  {
		//循环读取数据
		line, ok := <- tailfile.TailObj.Lines // chan tail.Line
		if !ok{
			logrus.Warn("tail file close reopen, filename:%s\n", tailfile.TailObj.Filename)
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
		msg.Topic = "web_log"
		msg.Value = sarama.StringEncoder(line.Text)
		// 将消息发送到通道中
		kafka.ToMsgChan(msg)

	}
}
func main() {
	// ini 使用教程https://ini.unknwon.io/docs
	var configObj = new(Config)
	// 0. 读取配置文件
	err := ini.MapTo(configObj, "G:\\goproject\\go\\logAgent\\confing\\confing.ini")
	if err != nil{
		logrus.Errorf("load config failed,err:%v", err)
	}
	fmt.Printf("configObj: %#v\n", configObj)
	// 1. 初始化连接kafka
	err = kafka.Init([]string{configObj.KafkaConfig.Address}, configObj.KafkaConfig.ChangeSize)
	if err != nil {
		logrus.Errorf("init kafka failed, err:%v", err)
		return
	}
	logrus.Info("init kafka success!")
	// 初始化etcd操作
	err = etcd.Init([]string{configObj.EtcdConfig.Address})
	if err != nil {
		logrus.Error("init etcd failed, err:%v", err)
		return
	}
	logrus.Info("init etcd file success!")

	// 拉取日志收集配置项的函数
	allEtcdConf, err:= etcd.GetConf(configObj.EtcdConfig.CollectKey)
	if err != nil {
		logrus.Error("get etcd conf failed, err:%v", err)
		return
	}
	fmt.Println("etcd config", allEtcdConf)
	// 2. 根据配置中的日志路径使用tail去收集日志
	err = tailfile.Init(configObj.CollectConfig.LogFilePath)
	if err != nil {
		logrus.Error("init tail file failed, err:%v", err)
		return
	}
	logrus.Info("init tail file success!")
	// 3. 把日志通过sarama发往kafka
	err = run()
	if err != nil {
		logrus.Error("run failed ,err :%v", err)
		return
	}
}