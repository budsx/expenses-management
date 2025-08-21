package payment

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/budsx/expenses-management/model"
)

type PaymentAPI struct {
	paymentProcessorURL string
	httpClient          *http.Client
}

func NewPaymentAPI(paymentProcessorURL string) *PaymentAPI {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	return &PaymentAPI{paymentProcessorURL: paymentProcessorURL, httpClient: client}
}

func (p *PaymentAPI) ProcessPayment(ctx context.Context, payment *model.PaymentProcessorModel) error {
	return nil
}

