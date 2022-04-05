package main

import (
	"bot/core"
	"bot/messaging"
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	amqp, err := messaging.Connect()
	if err != nil {
		return
	}
	defer amqp.Close()
	ch, err := messaging.OpenChannel()
	if err != nil {
		return
	}
	defer ch.Close()

	msgs := messaging.ReceiveMessageDeliveryChannel()

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var cm messaging.StockMessage
			json.Unmarshal(d.Body, &cm)
			message, err := core.GetStockQuote(cm.Message)
			fmt.Println(cm)
			stockMessage := &messaging.StockMessage{HubName: cm.HubName, ClientRemoteAddress: cm.ClientRemoteAddress, Message: message, RoomId: cm.RoomId}
			if err != nil {
				log.Fatal(err)
				return
			}

			messaging.SendMessage(stockMessage)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
