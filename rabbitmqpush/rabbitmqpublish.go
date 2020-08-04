package main

import (
	RabbitMQ "golearn/Rabbitmq"
	"strconv"
	"time"
)

func main() {
	kutengOne := RabbitMQ.NewRabbitMQTopic("exKutengTopic", "kuteng.topic.one", "amqp://sjd:sjd@192.168.1.171:5672")
	//kutengTwo := RabbitMQ.NewRabbitMQTopic("exKutengTopic", "kuteng.topic.two")
	//kutengThree := RabbitMQ.NewRabbitMQTopic("exKutengTopic", "kuteng.topic.three")
	go func() {
		t1 := time.Now()
		for i := 0; i <= 300000; i++ {
			kutengOne.PublishTopic("Hello kuteng topic one!" + strconv.Itoa(i))
			//kutengTwo.PublishTopic("Hello kuteng topic Two!" + strconv.Itoa(i))
			//time.Sleep(1 * time.Second)
			//fmt.Println(i)
		}
		println(time.Since(t1))
		kutengOne.Destory()
	}()
	/*go func() {
			for i := 0; i <= 1000000; i++ {
				//kutengOne.PublishTopic("Hello kuteng topic one!" + strconv.Itoa(i))
				kutengTwo.PublishTopic("Hello kuteng topic Two!" + strconv.Itoa(i))
				//time.Sleep(1 * time.Second)
				//fmt.Println(i)
			}

		}()

		go func() {
			for i := 0; i <= 1000000; i++ {
				//kutengOne.PublishTopic("Hello kuteng topic one!" + strconv.Itoa(i))
				kutengThree.PublishTopic("Hello kuteng topic three!" + strconv.Itoa(i))
				//time.Sleep(1 * time.Second)
				//fmt.Println(i)
			}

	    }()*/

	select {}
}
