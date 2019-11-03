package text

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type MyRabbitmq interface {
	Send(msg string)
	Rec() string
}
type RabbitStruct struct {
}

func rabbitmqerr(err error, msg string) {
	if err != nil {
		log.Fatal(err, msg)
	}
}
func (r RabbitStruct) Send(ququename, msg string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	rabbitmqerr(err, "连接失败")
	defer conn.Close()

	ch, err := conn.Channel()
	rabbitmqerr(err, "通道打开失败")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		ququename, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	rabbitmqerr(err, "声明队列失败")
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})
	rabbitmqerr(err, "发送信息失败")
	fmt.Println("发送的信息：", msg)
}
func (r RabbitStruct) Rec(ququename string) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	rabbitmqerr(err, "连接失败")
	defer conn.Close()

	ch, err := conn.Channel()
	rabbitmqerr(err, "通道打开失败")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		ququename, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	rabbitmqerr(err, "声明队列失败")

	for {
		msgs, ok, err := ch.Get(q.Name, false)
		rabbitmqerr(err, "get失败")
		if ok == true {
			log.Printf("收到的信息: %s", msgs.Body)
		} else {
			break
		}
	}
}
