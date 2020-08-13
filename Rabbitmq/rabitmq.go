package RabbitMQ

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

//rabbitMQ结构体
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	Mqurl   string
}

//创建结构体实例
func NewRabbitMQ(MQURL string) *RabbitMQ {
	return &RabbitMQ{Mqurl: MQURL}
}

//断开channel 和 connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

//错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

//话题模式
//创建RabbitMQ实例
func NewRabbitMQTopic(MQURL string) *RabbitMQ {
	//创建RabbitMQ实例
	rabbitmq := NewRabbitMQ(MQURL)
	var err error
	//获取connection
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "failed to connect rabbitmq!")
	//获取channel
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "failed to open a channel")
	return rabbitmq
}

//话题模式发送消息
func (r *RabbitMQ) PublishTopic(message string, exchangeName string, topickey string) {
	//1.尝试创建交换机
	err := r.channel.ExchangeDeclare(
		exchangeName,
		//要改成topic
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)

	r.failOnErr(err, "Failed to declare an excha"+
		"nge")

	//2.发送消息
	err = r.channel.Publish(
		exchangeName,
		//要设置
		topickey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}
func (r *RabbitMQ) RecieveTopic() {
	//1.试探性创建交换机
	err := r.channel.ExchangeDeclare(
		"SCCSIPSTATUS",
		//交换机类型
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exch"+
		"ange")
	//2.试探性创建队列，这里注意队列名称不要写
	q, err := r.channel.QueueDeclare(
		"", //随机生产队列名称
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare a queue")

	//绑定队列到 exchange 中
	err = r.channel.QueueBind(
		q.Name,
		//在pub/sub模式下，这里的key要为空
		"SCC.*.*",
		"SCCSIPSTATUS",
		false,
		nil)

	//消费消息
	messges, err := r.channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	forever := make(chan bool)

	go func() {
		for d := range messges {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	fmt.Println("退出请按 CTRL+C\n")
	<-forever
}
