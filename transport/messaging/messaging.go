package messaging

import (
	"context"
	"encoding/json"
	"log"

	"github.com/budsx/expenses-management/entity"
	"github.com/budsx/expenses-management/model"
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

func ProcessPaymentListener(service *service.ExpensesManagementService) func([]byte) error {
	return func(paymentResponse []byte) error {
		var payment *entity.PublishPaymentRequest
		err := json.Unmarshal(paymentResponse, &payment)
		if err != nil {
			return err
		}

		err = service.ProcessPayment(context.Background(), model.ApprovalRequest{
			ExpenseID:  payment.ExpenseID,
			ApproverID: payment.ApproverID,
			Notes:      payment.Notes,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
