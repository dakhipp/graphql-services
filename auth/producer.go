package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/dakhipp/graphql-services/auth/pb"
	"github.com/dakhipp/graphql-services/email"
	"github.com/dakhipp/graphql-services/text"
	"github.com/segmentio/kafka-go"
)

// Kafka is an interface which produces Kafka messages
type Kafka interface {
	RegisterEmail(ctx context.Context, args *pb.TriggerVerifyEmailRequest, code string) error
	ConfirmPhone(ctx context.Context, args *pb.TriggerVerifyPhoneRequest, code string) error
}

type kafkaProducer struct {
	emailProducer *kafka.Conn
	textProducer  *kafka.Conn
}

// NewKafkaProducer initializes kafka connections to all topics needed
func NewKafkaProducer(kafkaAddr string) (Kafka, error) {
	// create a connection on the email topic
	ep, err := kafka.DialLeader(context.Background(), "tcp", kafkaAddr, "email", 0)
	if err != nil {
		return nil, err
	}

	// create a connection on the text topic
	tp, err := kafka.DialLeader(context.Background(), "tcp", kafkaAddr, "text", 0)
	if err != nil {
		return nil, err
	}

	// attach both connections to the kafkaProducer struct
	return &kafkaProducer{ep, tp}, nil
}

// RegisterEmail takes a user and uses their info to send an email verification email
func (k *kafkaProducer) RegisterEmail(ctx context.Context, args *pb.TriggerVerifyEmailRequest, code string) error {
	// build kafka argument body
	a := email.ConfirmAccountArgs{
		FirstName:        args.FirstName,
		Email:            args.Email,
		VerificationCode: code,
	}

	// marshal email.ConfirmAccountArgs into byte array
	b, _ := json.Marshal(a)

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

	// log when message is sent
	fmt.Println("Message sent on topic \"email\" with key \"confirm-account\"")

	return nil
}

// ConfirmPhone takes a user and uses their info to send a phone verification text message
func (k *kafkaProducer) ConfirmPhone(ctx context.Context, args *pb.TriggerVerifyPhoneRequest, code string) error {
	// build kafka argument body
	a := text.ConfirmPhoneArgs{
		ToPhone:          args.Phone,
		VerificationCode: code,
	}

	// marshal text.ConfirmPhoneArgs into byte array
	b, _ := json.Marshal(a)

	// create kafka message the key identifies which template to use and the value contains arguments needed to compile the template
	m := kafka.Message{
		Key:   []byte("confirm-phone"),
		Value: b,
	}

	// write message to kafka topic
	_, err := k.textProducer.WriteMessages(m)
	if err != nil {
		fmt.Println(err)
	}

	// log when message is sent
	fmt.Println("Message sent on topic \"text\" with key \"confirm-phone\"")

	return nil
}
