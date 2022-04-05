package messager

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/ezrod12/chat/models"
	"github.com/ezrod12/chat/settings"
	"github.com/streadway/amqp"
)

var Conn *amqp.Connection
var Channel *amqp.Channel
var GetStockQueue, SendStockQueue string

func Connect(cfg *settings.Config) (*amqp.Connection, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
		cfg.RabbitMQ.User,
		cfg.RabbitMQ.Pass,
		cfg.RabbitMQ.Host,
		cfg.RabbitMQ.Port))

	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
		return nil, err
	}

	Conn = conn

	GetStockQueue = cfg.RabbitMQ.GetStockQueue
	SendStockQueue = cfg.RabbitMQ.SendStockQueue

	return conn, nil
}

func OpenChannel() (*amqp.Channel, error) {
	ch, err := Conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open channel")
		return nil, err
	}

	Channel = ch

	return ch, nil
}

func SendMessage(message *models.StockMessage) {
	q, err := Channel.QueueDeclare(
		GetStockQueue, // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	json, err := json.Marshal(message)
	if err != nil {
		failOnError(err, "Failed to parse body message")
	}

	err = Channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        json,
		})
	failOnError(err, "Failed to publish a message")

	fmt.Printf("Message sent: %s\n", json)
}

func ReceiveMessageDeliveryChannel() <-chan amqp.Delivery {
	q, err := Channel.QueueDeclare(
		SendStockQueue, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	return msgs
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
