package kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

/*
@author RandySun
@create 2021-09-11-20:58
*/
var (
	client sarama.SyncProducer
	msgChan chan *sarama.ProducerMessage
)
func  Init(address [] string, chanSize int64) (err error) {
	// 1. 生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 随机选出一个partition
	config.Producer.Return.Successes = true // 成功交付的消息将在success channel返回
	//  2. 连接kafka
	client, err = sarama.NewSyncProducer(address, config)
	if err != nil{
		fmt.Println("producer closed  err: ", err)
		return
	}
	//fmt.Printf(client, "")
	// 初始化MsgChan
	msgChan = make(chan *sarama.ProducerMessage, chanSize)
	go sendMsg() // 向kafka发送消息
	return
}

// 从Msgchan中读取Msg，发送给Kafka

func sendMsg(){
	for {
		select {
		case msg := <-msgChan:
			pid, offset, err := client.SendMessage(msg)
			if err != nil{
				fmt.Println("sned msg failed err: ", err)
			}
			logrus.Infof("send msg to kafka success. pid:%v offset:%v, msg: %v", pid, offset, msg)
		}
	}
}

// 定义一个函数向外暴露msgChan
func ToMsgChan(msg *sarama.ProducerMessage) {
	msgChan <-msg
}