package rabbitmq

import (
	"encoding/json"

	"github.com/budsx/expenses-management/entity"
	"github.com/budsx/expenses-management/util/rabbitmq"
)

type RabbitMQClient struct {
	client                *rabbitmq.RabbitMQClient
	topicPaymentProcessor string
}

func NewRabbitClient(client *rabbitmq.RabbitMQClient, topicPaymentProcessor string) *RabbitMQClient {
	return &RabbitMQClient{
		client:                client,
		topicPaymentProcessor: topicPaymentProcessor,
	}
}

func (c *RabbitMQClient) PublishPayment(payment *entity.PublishPaymentRequest) error {
	jsonData, err := json.Marshal(payment)
	if err != nil {
		return err
	}
	return c.client.Publish(c.topicPaymentProcessor, jsonData)
}

func (c *RabbitMQClient) GetClient() *rabbitmq.RabbitMQClient {
	if c.client != nil {
		return c.client
	}
	return nil
}
