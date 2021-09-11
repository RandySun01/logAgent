package main

/*
@author RandySun
@create 2021-09-11-11:23
*/
import (
	"fmt"
	"github.com/Shopify/sarama"
)
// 基于sarama第三⽅库开发的kafka client
func main() {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 随机选出一个partition
	config.Producer.Return.Successes = true // 成功交付的消息将在success channel返回
	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "web_log"
	msg.Value = sarama.StringEncoder("this is a randy sun log")
	// 连接kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil{
		fmt.Println("producer closed  err: ", err)
	}
	defer client.Close()
	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil{
		fmt.Println("sned msg failed err: ", err)
	}
	fmt.Printf("pid:%v offset: %v\n", pid, offset)

	fmt.Println(config)
}