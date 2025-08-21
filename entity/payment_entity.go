package entity

type PaymentProcessorRequest struct {
	AmountIDR  int    `json:"amount_idr"`
	ExternalID string `json:"external_id"`
}

type PaymentProcessorResponse struct {
	Data struct {
		ID         string `json:"id"`
		ExternalID string `json:"external_id"`
		Status     string `json:"status"`
	} `json:"data"`
	Message string `json:"message"`
}