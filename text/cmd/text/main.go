package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dakhipp/graphql-services/text"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/segmentio/kafka-go"
)

// funcMap is a map that maps kafka message keys to functions used to send out different text messages
var funcMap = map[string]interface{}{
	"confirm-phone": text.SendConfirmPhone,
}

type envConfig struct {
	KafkaBrokers string `envconfig:"KAFKA_BROKERS"`
}

func main() {
	// declare and attempt to cast config struct
	var cfg envConfig
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// set up Kafka configuration
	c := kafka.ReaderConfig{
		Brokers:  strings.Split(cfg.KafkaBrokers, ","),
		GroupID:  "text-consumer-group",
		Topic:    "text",
		MinBytes: 10e3,            // 10KB
		MaxBytes: 10e6,            // 10MB
		MaxWait:  1 * time.Second, // Maximum amount of time to wait for new data to come when fetching batches of messages from kafka
	}

	// create new kafka reader
	r := kafka.NewReader(c)
	fmt.Println("Kafka consumer listening on topic \"text\"")

	// listen for messages with Kafka config
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("there was an error reading the message %v\n", err)
			break
		}

		key := string(m.Key)

		fmt.Printf("new message read with key '%s'\n", key)

		if _, ok := funcMap[key]; ok {
			funcMap[key].(func([]byte))(m.Value)
		} else {
			fmt.Printf("the template key '%v' was not found\n", key)
		}
	}

	r.Close()
}
