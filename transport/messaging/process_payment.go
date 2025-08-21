package messaging

import (
	"context"
	"encoding/json"

	"github.com/budsx/expenses-management/entity"
	"github.com/budsx/expenses-management/model"
	"github.com/budsx/expenses-management/service"
)

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
			Status:     payment.Status,
		})
		if err != nil {
			return err
		}

		return nil
	}
}
