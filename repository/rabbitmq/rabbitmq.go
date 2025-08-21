package rabbitmq

import (
	"encoding/json"

	"github.com/budsx/expenses-management/entity"
	"github.com/budsx/expenses-management/util/rabbitmq"
)

type RabbitMQClient struct {
	client *rabbitmq.RabbitMQClient
}

func NewRabbitClient(client *rabbitmq.RabbitMQClient) *RabbitMQClient {
	return &RabbitMQClient{
		client: client,
	}
}

func (c *RabbitMQClient) PublishPayment(topic string, payment *entity.PublishPaymentRequest) error {
	jsonData, err := json.Marshal(payment)
	if err != nil {
		return err
	}
	return c.client.Publish(topic, jsonData)
}

func (c *RabbitMQClient) GetClient() *rabbitmq.RabbitMQClient {
	if c.client != nil {
		return c.client
	}
	return nil
}
