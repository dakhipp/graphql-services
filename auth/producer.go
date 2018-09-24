package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dakhipp/graphql-services/email"
	"github.com/segmentio/kafka-go"
)

// Kafka is an interface which produces Kafka messages
type Kafka interface {
	RegisterEmail(ctx context.Context, args User) error
}

type kafkaProducer struct {
	emailProducer *kafka.Conn
}

// NewKafkaProducer initializes a new database connection for the repository
func NewKafkaProducer(kafkaAddr string) (Kafka, error) {
	ep, err := kafka.DialLeader(context.Background(), "tcp", kafkaAddr, "email", 0)
	if err != nil {
		return nil, err
	}

	return &kafkaProducer{ep}, nil
}

// RegisterEmail takes a user and uses their info to
func (k *kafkaProducer) RegisterEmail(ctx context.Context, u User) error {
	// build kafka argument body
	args := email.ConfirmAccountArgs{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Email:     u.Email,
	}

	// marshal email.ConfirmAccountArgs into byte array
	b, _ := json.Marshal(args)

	// create kafka message the key identifies which template to use and the value contains arguments needed to compile the template
	m := kafka.Message{
		Key:   []byte("confirm-account"),
		Value: b,
	}

	// write message to kafka topic
	_, err := k.emailProducer.WriteMessages(m)
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
