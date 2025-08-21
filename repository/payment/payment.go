package payment

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/budsx/expenses-management/entity"
)

type paymentProcessor struct {
	paymentProcessorURL string
	httpClient          *http.Client
}

var (
	once   sync.Once
	client *paymentProcessor
)

func NewPaymentProcessor(paymentProcessorURL string) *paymentProcessor {
	once.Do(func() {
		client = &paymentProcessor{
			paymentProcessorURL: paymentProcessorURL,
			httpClient: &http.Client{
				Timeout: time.Second * 30,
				Transport: &http.Transport{
					Dial: (&net.Dialer{
						Timeout:   time.Second * 30,
						KeepAlive: time.Second * 30,
					}).Dial,
					TLSHandshakeTimeout: time.Second * 10,
				},
			},
		}
	})
	return client
}

func (p *paymentProcessor) GetClient() *http.Client {
	return p.httpClient
}

func (p *paymentProcessor) ProcessPayment(ctx context.Context, payment *entity.PaymentProcessorRequest) (*entity.PaymentProcessorResponse, error) {
	client := p.GetClient()
	requestBody, err := json.Marshal(payment)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, p.paymentProcessorURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to process payment: %s", response.Status)
	}

	var responseBody *entity.PaymentProcessorResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}
