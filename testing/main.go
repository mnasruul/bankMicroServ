package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/Shopify/sarama"
)

var (
	brokerList        = kingpin.Flag("brokerList", "List of brokers to connect").Default("192.168.223.109:9092").Strings()
	topic             = kingpin.Flag("topic", "Topic name").Default("xbanku-transactions-t3").String()
	partition         = kingpin.Flag("partition", "Partition number").Default("0").String()
	offsetType        = kingpin.Flag("offsetType", "Offset Type (OffsetNewest | OffsetOldest)").Default("-1").Int()
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
)

func main() {
	kingpin.Parse()
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := *brokerList
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			log.Panic(err)
		}
	}()
	consumer, err := master.ConsumePartition(*topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Panic(err)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})
	var log interface{}
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				*messageCountStart++

				msgVal := msg.Value
				// fmt.Println(msgVal)
				if err = json.Unmarshal(msgVal, &log); err != nil {
					fmt.Printf("Failed parsing: %s", err)
				} else {
					logMap := log.(map[string]interface{})
					logType := logMap["Type"]
					// fmt.Println(logType)
					fmt.Printf("Processing %s:\n%s\n", logType, string(msgVal))
					fmt.Println("Received messages", string(msg.Key), string(msg.Value))
				}

			case <-signals:
				fmt.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	fmt.Println("Processed", *messageCountStart, "messages")
}
