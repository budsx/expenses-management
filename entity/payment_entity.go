package entity

type PaymentProcessorRequest struct {
	AmountIDR  int64  `json:"amount_idr"`
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

type PublishPaymentRequest struct {
	ExpenseID  int64  `json:"expense_id"`
	ApproverID int64  `json:"approver_id"`
	Notes      string `json:"notes"`
	Status     int32  `json:"status"`
}
