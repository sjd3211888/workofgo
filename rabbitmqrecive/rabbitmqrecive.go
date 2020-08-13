package main

import RabbitMQ "golearn/Rabbitmq"

func main() {
	kutengOne := RabbitMQ.NewRabbitMQTopic("amqp://sjd:sjd@39.107.237.49:5672/admin")
	kutengOne.RecieveTopic()
}
