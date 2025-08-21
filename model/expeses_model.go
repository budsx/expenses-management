package model

type CreateExpenseRequest struct {
	AmountIDR   float64 `json:"amount_idr" validate:"required,gt=0"`
	Description string  `json:"description" validate:"required"`
	ReceiptURL  string  `json:"receipt_url"`
}

type UpdateExpenseStatusRequest struct {
	Status string `json:"status" validate:"required"`
	Notes  string `json:"notes"`
}

type ExpenseResponse struct {
	ID           int64   `json:"id"`
	UserID       int64   `json:"user_id"`
	AmountIDR    float64 `json:"amount_idr"`
	Description  string  `json:"description"`
	ReceiptURL   string  `json:"receipt_url"`
	Status       string  `json:"status"`
	AutoApproved bool    `json:"auto_approved"`
}

type ExpenseListResponse struct {
	Expenses   []ExpenseResponse `json:"expenses"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

type ExpenseListQuery struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Status   string `query:"status"`
	UserID   int64  `query:"user_id"`
}

type ApprovalResponse struct {
	Message string `json:"message"`
}
