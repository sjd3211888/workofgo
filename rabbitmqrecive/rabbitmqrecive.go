package main

import RabbitMQ "golearn/Rabbitmq"

func main() {
	kutengOne := RabbitMQ.NewRabbitMQTopic("exKutengTopic", "#", "amqp://sjd:sjd@192.168.1.171:5672")
	kutengOne.RecieveTopic()
}
