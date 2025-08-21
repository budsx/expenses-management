package messaging

import (
	"log"

	_interface "github.com/budsx/expenses-management/repository/interface"
	"github.com/budsx/expenses-management/service"
)

func NewTransportListener(service *service.ExpensesManagementService, rabbitmqClient _interface.RabbitMQClient, exchangeName string, queueName string) {
	client := rabbitmqClient.GetClient()

	if err := client.Subscribe(exchangeName, queueName, ProcessPaymentListener(service)); err != nil {
		log.Printf("Failed to subscribe to %s: %v", queueName, err)
		return
	}
}
