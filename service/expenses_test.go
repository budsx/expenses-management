package service

import (
	"context"
	"errors"
	"testing"

	"github.com/budsx/expenses-management/entity"
	"github.com/budsx/expenses-management/model"
	"github.com/budsx/expenses-management/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestExpensesService_CreateExpense(t *testing.T) {
	tests := []struct {
		name    string
		request model.CreateExpenseRequest
		userCtx model.User
		mock    func(server *TestService)
		want    *model.ExpenseResponse
		wantErr bool
	}{
		{
			name: "success - auto approved",
			request: model.CreateExpenseRequest{
				AmountIDR:   100000,
				Description: "Test Expense",
				ReceiptURL:  "https://example.com/receipt.jpg",
			},
			userCtx: model.User{
				ID:    1,
				Email: "test@example.com",
				Role:  int(util.USER_ROLE_EMPLOYEE),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					WriteExpense(gomock.Any(), gomock.Any()).
					Return(int64(123), nil).
					Times(1)

				server.MockRepo.EXPECT().
					WriteAuditLog(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				server.MockRabbitMQ.EXPECT().
					PublishPayment(gomock.Any()).
					Return(nil).
					AnyTimes()
			},
			want: &model.ExpenseResponse{
				ID:           123,
				UserID:       1,
				AmountIDR:    100000,
				Description:  "Test Expense",
				ReceiptURL:   "https://example.com/receipt.jpg",
				Status:       util.GetExpenseStatusString(util.EXPENSE_PENDING),
				AutoApproved: true,
			},
			wantErr: false,
		},
		{
			name: "success - manual approval required",
			request: model.CreateExpenseRequest{
				AmountIDR:   2000000,
				Description: "Large Expense",
				ReceiptURL:  "https://example.com/receipt2.jpg",
			},
			userCtx: model.User{
				ID:    2,
				Email: "user2@example.com",
				Role:  int(util.USER_ROLE_EMPLOYEE),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					WriteExpense(gomock.Any(), gomock.Any()).
					Return(int64(456), nil).
					Times(1)

				server.MockRepo.EXPECT().
					WriteAuditLog(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			want: &model.ExpenseResponse{
				ID:           456,
				UserID:       2,
				AmountIDR:    2000000,
				Description:  "Large Expense",
				ReceiptURL:   "https://example.com/receipt2.jpg",
				Status:       util.GetExpenseStatusString(util.EXPENSE_PENDING),
				AutoApproved: false,
			},
			wantErr: false,
		},
		{
			name: "write expense error",
			request: model.CreateExpenseRequest{
				AmountIDR:   50000,
				Description: "Failed Expense",
				ReceiptURL:  "https://example.com/receipt3.jpg",
			},
			userCtx: model.User{
				ID:    3,
				Email: "user3@example.com",
				Role:  int(util.USER_ROLE_EMPLOYEE),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					WriteExpense(gomock.Any(), gomock.Any()).
					Return(int64(0), errors.New("database error")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "audit log error - should not fail",
			request: model.CreateExpenseRequest{
				AmountIDR:   75000,
				Description: "Audit Log Error",
				ReceiptURL:  "https://example.com/receipt4.jpg",
			},
			userCtx: model.User{
				ID:    4,
				Email: "user4@example.com",
				Role:  int(util.USER_ROLE_EMPLOYEE),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					WriteExpense(gomock.Any(), gomock.Any()).
					Return(int64(789), nil).
					Times(1)

				server.MockRepo.EXPECT().
					WriteAuditLog(gomock.Any(), gomock.Any()).
					Return(errors.New("audit log error")).
					Times(1)

				server.MockRabbitMQ.EXPECT().
					PublishPayment(gomock.Any()).
					Return(nil).
					AnyTimes()
			},
			want: &model.ExpenseResponse{
				ID:           789,
				UserID:       4,
				AmountIDR:    75000,
				Description:  "Audit Log Error",
				ReceiptURL:   "https://example.com/receipt4.jpg",
				Status:       util.GetExpenseStatusString(util.EXPENSE_PENDING),
				AutoApproved: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t)
			defer server.MockCtrl.Finish()

			ctx := context.Background()
			ctx = context.WithValue(ctx, "user_id", tt.userCtx.ID)
			ctx = context.WithValue(ctx, "user_email", tt.userCtx.Email)
			ctx = context.WithValue(ctx, "user_role", tt.userCtx.Role)

			tt.mock(server)

			got, err := server.Service.CreateExpense(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.UserID, got.UserID)
			assert.Equal(t, tt.want.AmountIDR, got.AmountIDR)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.Equal(t, tt.want.ReceiptURL, got.ReceiptURL)
			assert.Equal(t, tt.want.AutoApproved, got.AutoApproved)
		})
	}
}

func TestExpensesService_CreateExpense_InvalidContext(t *testing.T) {
	server := NewTestServer(t)
	defer server.MockCtrl.Finish()
	ctx := context.Background()

	_, err := server.Service.CreateExpense(ctx, model.CreateExpenseRequest{
		AmountIDR:   100000,
		Description: "Test Expense",
		ReceiptURL:  "https://example.com/receipt.jpg",
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user info")
}

func TestExpensesService_GetExpenses(t *testing.T) {
	tests := []struct {
		name    string
		query   model.ExpenseListQuery
		userCtx model.User
		mock    func(server *TestService)
		want    *model.ExpenseListResponse
		wantErr bool
	}{
		{
			name: "success - manager gets all expenses",
			query: model.ExpenseListQuery{
				Page:     1,
				PageSize: 10,
				Status:   0,
				UserID:   0,
			},
			userCtx: model.User{
				ID:    1,
				Email: "manager@example.com",
				Role:  int(util.USER_ROLE_MANAGER),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpensesWithPagination(gomock.Any(), gomock.Any()).
					Return([]*entity.Expense{
						{
							ID:           1,
							UserID:       2,
							AmountIDR:    100000,
							Description:  "Test Expense 1",
							ReceiptURL:   "https://example.com/receipt1.jpg",
							Status:       int32(util.EXPENSE_PENDING),
							AutoApproved: true,
						},
						{
							ID:           2,
							UserID:       3,
							AmountIDR:    200000,
							Description:  "Test Expense 2",
							ReceiptURL:   "https://example.com/receipt2.jpg",
							Status:       int32(util.EXPENSE_APPROVED),
							AutoApproved: false,
						},
					}, int64(2), nil).
					Times(1)
			},
			want: &model.ExpenseListResponse{
				Expenses: []model.ExpenseResponse{
					{
						ID:           1,
						UserID:       2,
						AmountIDR:    100000,
						Description:  "Test Expense 1",
						ReceiptURL:   "https://example.com/receipt1.jpg",
						Status:       util.GetExpenseStatusString(util.EXPENSE_PENDING),
						AutoApproved: true,
					},
					{
						ID:           2,
						UserID:       3,
						AmountIDR:    200000,
						Description:  "Test Expense 2",
						ReceiptURL:   "https://example.com/receipt2.jpg",
						Status:       util.GetExpenseStatusString(util.EXPENSE_APPROVED),
						AutoApproved: false,
					},
				},
				Total:    2,
				Page:     1,
				PageSize: 10,
			},
			wantErr: false,
		},
		{
			name: "success - employee gets own expenses only",
			query: model.ExpenseListQuery{
				Page:     1,
				PageSize: 5,
				Status:   0,
				UserID:   0,
			},
			userCtx: model.User{
				ID:    2,
				Email: "employee@example.com",
				Role:  int(util.USER_ROLE_EMPLOYEE),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpensesWithPagination(gomock.Any(), gomock.Any()).
					Return([]*entity.Expense{
						{
							ID:           1,
							UserID:       2,
							AmountIDR:    50000,
							Description:  "Employee Expense",
							ReceiptURL:   "https://example.com/receipt.jpg",
							Status:       int32(util.EXPENSE_PENDING),
							AutoApproved: true,
						},
					}, int64(1), nil).
					Times(1)
			},
			want: &model.ExpenseListResponse{
				Expenses: []model.ExpenseResponse{
					{
						ID:           1,
						UserID:       2,
						AmountIDR:    50000,
						Description:  "Employee Expense",
						ReceiptURL:   "https://example.com/receipt.jpg",
						Status:       util.GetExpenseStatusString(util.EXPENSE_PENDING),
						AutoApproved: true,
					},
				},
				Total:    1,
				Page:     1,
				PageSize: 5,
			},
			wantErr: false,
		},
		{
			name: "database error",
			query: model.ExpenseListQuery{
				Page:     1,
				PageSize: 10,
			},
			userCtx: model.User{
				ID:    1,
				Email: "manager@example.com",
				Role:  int(util.USER_ROLE_MANAGER),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpensesWithPagination(gomock.Any(), gomock.Any()).
					Return(nil, int64(0), errors.New("database error")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t)
			defer server.MockCtrl.Finish()

			ctx := context.Background()
			ctx = context.WithValue(ctx, "user_id", tt.userCtx.ID)
			ctx = context.WithValue(ctx, "user_email", tt.userCtx.Email)
			ctx = context.WithValue(ctx, "user_role", tt.userCtx.Role)

			tt.mock(server)

			got, err := server.Service.GetExpenses(ctx, tt.query)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.Total, got.Total)
			assert.Equal(t, tt.want.Page, got.Page)
			assert.Equal(t, tt.want.PageSize, got.PageSize)
			assert.Equal(t, len(tt.want.Expenses), len(got.Expenses))
		})
	}
}

func TestExpensesService_GetExpenses_InvalidContext(t *testing.T) {
	server := NewTestServer(t)
	defer server.MockCtrl.Finish()
	ctx := context.Background()

	_, err := server.Service.GetExpenses(ctx, model.ExpenseListQuery{
		Page:     1,
		PageSize: 10,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get user info")
}

func TestExpensesService_GetExpenseByID(t *testing.T) {
	tests := []struct {
		name      string
		expenseID int64
		mock      func(server *TestService)
		want      *model.ExpenseResponse
		wantErr   bool
	}{
		{
			name:      "success",
			expenseID: 123,
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpenseByID(gomock.Any(), int64(123)).
					Return(&entity.Expense{
						ID:           123,
						UserID:       1,
						AmountIDR:    150000,
						Description:  "Test Expense",
						ReceiptURL:   "https://example.com/receipt.jpg",
						Status:       int32(util.EXPENSE_APPROVED),
						AutoApproved: false,
					}, nil).
					Times(1)
			},
			want: &model.ExpenseResponse{
				ID:           123,
				UserID:       1,
				AmountIDR:    150000,
				Description:  "Test Expense",
				ReceiptURL:   "https://example.com/receipt.jpg",
				Status:       util.GetExpenseStatusString(util.EXPENSE_APPROVED),
				AutoApproved: false,
			},
			wantErr: false,
		},
		{
			name:      "expense not found",
			expenseID: 999,
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpenseByID(gomock.Any(), int64(999)).
					Return(nil, errors.New("expense not found")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:      "database error",
			expenseID: 123,
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpenseByID(gomock.Any(), int64(123)).
					Return(nil, errors.New("database connection failed")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t)
			defer server.MockCtrl.Finish()

			ctx := context.Background()
			tt.mock(server)

			got, err := server.Service.GetExpenseByID(ctx, tt.expenseID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.ID, got.ID)
			assert.Equal(t, tt.want.UserID, got.UserID)
			assert.Equal(t, tt.want.AmountIDR, got.AmountIDR)
			assert.Equal(t, tt.want.Description, got.Description)
			assert.Equal(t, tt.want.ReceiptURL, got.ReceiptURL)
			assert.Equal(t, tt.want.Status, got.Status)
			assert.Equal(t, tt.want.AutoApproved, got.AutoApproved)
		})
	}
}

func TestExpensesService_ApproveExpense(t *testing.T) {
	tests := []struct {
		name    string
		request model.ApprovalRequest
		userCtx model.User
		mock    func(server *TestService)
		want    *model.ApprovalResponse
		wantErr bool
		errMsg  string
	}{
		{
			name: "success - manager approves expense",
			request: model.ApprovalRequest{
				ExpenseID:  123,
				ApproverID: 2,
				Notes:      "Approved by manager",
				Status:     int32(util.APPROVAL_APPROVED),
			},
			userCtx: model.User{
				ID:    2,
				Email: "manager@example.com",
				Role:  int(util.USER_ROLE_MANAGER),
			},
			mock: func(server *TestService) {
				server.MockRabbitMQ.EXPECT().
					PublishPayment(gomock.Any()).
					Return(nil).
					AnyTimes()
			},
			want: &model.ApprovalResponse{
				Message: "Expense 123 approved",
			},
			wantErr: false,
		},
		{
			name: "failure - not a manager",
			request: model.ApprovalRequest{
				ExpenseID:  123,
				ApproverID: 1,
				Notes:      "Should fail",
				Status:     int32(util.APPROVAL_APPROVED),
			},
			userCtx: model.User{
				ID:    1,
				Email: "employee@example.com",
				Role:  int(util.USER_ROLE_EMPLOYEE),
			},
			mock:    func(server *TestService) {},
			want:    nil,
			wantErr: true,
			errMsg:  "user is not a manager",
		},
		{
			name: "failure - invalid context",
			request: model.ApprovalRequest{
				ExpenseID: 123,
			},
			userCtx: model.User{}, // Empty context
			mock:    func(server *TestService) {},
			want:    nil,
			wantErr: true,
			errMsg:  "failed to get user info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t)
			defer server.MockCtrl.Finish()

			var ctx context.Context = context.Background()
			if tt.userCtx.ID != 0 {
				ctx = context.WithValue(ctx, "user_id", tt.userCtx.ID)
				ctx = context.WithValue(ctx, "user_email", tt.userCtx.Email)
				ctx = context.WithValue(ctx, "user_role", tt.userCtx.Role)
			}

			tt.mock(server)

			got, err := server.Service.ApproveExpense(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.Message, got.Message)
		})
	}
}

func TestExpensesService_RejectExpense(t *testing.T) {
	tests := []struct {
		name    string
		request model.ApprovalRequest
		userCtx model.User
		mock    func(server *TestService)
		want    *model.ApprovalResponse
		wantErr bool
		errMsg  string
	}{
		{
			name: "success - manager rejects expense",
			request: model.ApprovalRequest{
				ExpenseID:  123,
				ApproverID: 2,
				Notes:      "Insufficient documentation",
				Status:     int32(util.EXPENSE_REJECTED),
			},
			userCtx: model.User{
				ID:    2,
				Email: "manager@example.com",
				Role:  int(util.USER_ROLE_MANAGER),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					ApprovalExpense(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				server.MockRepo.EXPECT().
					WriteAuditLog(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			want: &model.ApprovalResponse{
				Message: "Expense 123 successfully rejected",
			},
			wantErr: false,
		},
		{
			name: "success - audit log error should not fail",
			request: model.ApprovalRequest{
				ExpenseID:  123,
				ApproverID: 2,
				Notes:      "Rejected",
				Status:     int32(util.EXPENSE_REJECTED),
			},
			userCtx: model.User{
				ID:    2,
				Email: "manager@example.com",
				Role:  int(util.USER_ROLE_MANAGER),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					ApprovalExpense(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				server.MockRepo.EXPECT().
					WriteAuditLog(gomock.Any(), gomock.Any()).
					Return(errors.New("audit log error")).
					Times(1)
			},
			want: &model.ApprovalResponse{
				Message: "Expense 123 successfully rejected",
			},
			wantErr: false,
		},
		{
			name: "failure - not a manager",
			request: model.ApprovalRequest{
				ExpenseID: 123,
			},
			userCtx: model.User{
				ID:    1,
				Email: "employee@example.com",
				Role:  int(util.USER_ROLE_EMPLOYEE),
			},
			mock:    func(server *TestService) {},
			want:    nil,
			wantErr: true,
			errMsg:  "user is not a manager",
		},
		{
			name: "failure - approval expense error",
			request: model.ApprovalRequest{
				ExpenseID:  123,
				ApproverID: 2,
				Notes:      "Should fail",
				Status:     int32(util.EXPENSE_REJECTED),
			},
			userCtx: model.User{
				ID:    2,
				Email: "manager@example.com",
				Role:  int(util.USER_ROLE_MANAGER),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					ApprovalExpense(gomock.Any(), gomock.Any()).
					Return(errors.New("database error")).
					Times(1)
			},
			want:    nil,
			wantErr: true,
			errMsg:  "failed to reject expense",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t)
			defer server.MockCtrl.Finish()

			var ctx context.Context = context.Background()
			if tt.userCtx.ID != 0 {
				ctx = context.WithValue(ctx, "user_id", tt.userCtx.ID)
				ctx = context.WithValue(ctx, "user_email", tt.userCtx.Email)
				ctx = context.WithValue(ctx, "user_role", tt.userCtx.Role)
			}

			tt.mock(server)

			got, err := server.Service.RejectExpense(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want.Message, got.Message)
		})
	}
}

func TestExpensesService_HealthCheck(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(server *TestService)
		wantErr bool
	}{
		{
			name: "success - database is healthy",
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					PingContext(gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "failure - database connection failed",
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					PingContext(gomock.Any()).
					Return(errors.New("connection refused")).
					Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServer(t)
			defer server.MockCtrl.Finish()

			ctx := context.Background()
			tt.mock(server)

			err := server.Service.HealthCheck(ctx)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestExpensesService_ProcessPayment(t *testing.T) {
	tests := []struct {
		name    string
		request model.ApprovalRequest
		mock    func(server *TestService)
		wantErr bool
		errMsg  string
	}{
		{
			name: "success - process payment for pending expense",
			request: model.ApprovalRequest{
				ExpenseID:  123,
				ApproverID: 2,
				Notes:      "Approved for payment",
				Status:     int32(util.EXPENSE_APPROVED),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpenseByID(gomock.Any(), int64(123)).
					Return(&entity.Expense{
						ID:           123,
						UserID:       1,
						AmountIDR:    150000,
						Description:  "Test Expense",
						ReceiptURL:   "https://example.com/receipt.jpg",
						Status:       int32(util.EXPENSE_PENDING),
						AutoApproved: false,
					}, nil).
					Times(1)

				server.MockRepo.EXPECT().
					ApprovalExpense(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				server.MockPaymentProcessor.EXPECT().
					ProcessPayment(gomock.Any(), gomock.Any()).
					Return(&entity.PaymentProcessorResponse{
						Data: struct {
							ID         string `json:"id"`
							ExternalID string `json:"external_id"`
							Status     string `json:"status"`
						}{
							ID:         "TXN123",
							ExternalID: "EXT123",
							Status:     "SUCCESS",
						},
						Message: "Payment processed successfully",
					}, nil).
					Times(1)

				server.MockRepo.EXPECT().
					WriteAuditLog(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "success - process payment for auto-approved expense",
			request: model.ApprovalRequest{
				ExpenseID:  124,
				ApproverID: 0,
				Notes:      "Auto approved payment",
				Status:     int32(util.EXPENSE_AUTO_APPROVED),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpenseByID(gomock.Any(), int64(124)).
					Return(&entity.Expense{
						ID:           124,
						UserID:       1,
						AmountIDR:    75000,
						Description:  "Auto Approved Expense",
						ReceiptURL:   "https://example.com/receipt.jpg",
						Status:       int32(util.EXPENSE_AUTO_APPROVED),
						AutoApproved: true,
					}, nil).
					Times(1)

				server.MockRepo.EXPECT().
					ApprovalExpense(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				server.MockPaymentProcessor.EXPECT().
					ProcessPayment(gomock.Any(), gomock.Any()).
					Return(&entity.PaymentProcessorResponse{
						Data: struct {
							ID         string `json:"id"`
							ExternalID string `json:"external_id"`
							Status     string `json:"status"`
						}{
							ID:         "TXN124",
							ExternalID: "EXT124",
							Status:     "SUCCESS",
						},
						Message: "Payment processed successfully",
					}, nil).
					Times(1)

				server.MockRepo.EXPECT().
					WriteAuditLog(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
		},
		{
			name: "failure - expense not found",
			request: model.ApprovalRequest{
				ExpenseID: 999,
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpenseByID(gomock.Any(), int64(999)).
					Return(nil, errors.New("expense not found")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to get expense",
		},
		{
			name: "failure - expense already processed",
			request: model.ApprovalRequest{
				ExpenseID: 123,
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpenseByID(gomock.Any(), int64(123)).
					Return(&entity.Expense{
						ID:           123,
						UserID:       1,
						AmountIDR:    150000,
						Description:  "Already Processed Expense",
						ReceiptURL:   "https://example.com/receipt.jpg",
						Status:       int32(util.EXPENSE_APPROVED), // Already approved
						AutoApproved: false,
					}, nil).
					Times(1)
			},
			wantErr: true,
			errMsg:  "expense is not pending",
		},
		{
			name: "failure - approval expense error",
			request: model.ApprovalRequest{
				ExpenseID:  123,
				ApproverID: 2,
				Notes:      "Should fail approval",
				Status:     int32(util.EXPENSE_APPROVED),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpenseByID(gomock.Any(), int64(123)).
					Return(&entity.Expense{
						ID:           123,
						UserID:       1,
						AmountIDR:    150000,
						Description:  "Test Expense",
						Status:       int32(util.EXPENSE_PENDING),
						AutoApproved: false,
					}, nil).
					Times(1)

				server.MockRepo.EXPECT().
					ApprovalExpense(gomock.Any(), gomock.Any()).
					Return(errors.New("approval failed")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to approve expense",
		},
		{
			name: "failure - payment processor error",
			request: model.ApprovalRequest{
				ExpenseID:  123,
				ApproverID: 2,
				Notes:      "Payment should fail",
				Status:     int32(util.EXPENSE_APPROVED),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpenseByID(gomock.Any(), int64(123)).
					Return(&entity.Expense{
						ID:           123,
						UserID:       1,
						AmountIDR:    150000,
						Description:  "Test Expense",
						Status:       int32(util.EXPENSE_PENDING),
						AutoApproved: false,
					}, nil).
					Times(1)

				server.MockRepo.EXPECT().
					ApprovalExpense(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				server.MockPaymentProcessor.EXPECT().
					ProcessPayment(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("payment processor unavailable")).
					Times(1)
			},
			wantErr: true,
			errMsg:  "failed to process payment",
		},
		{
			name: "success - audit log error should not fail",
			request: model.ApprovalRequest{
				ExpenseID:  123,
				ApproverID: 2,
				Notes:      "Audit log will fail",
				Status:     int32(util.EXPENSE_APPROVED),
			},
			mock: func(server *TestService) {
				server.MockRepo.EXPECT().
					GetExpenseByID(gomock.Any(), int64(123)).
					Return(&entity.Expense{
						ID:           123,
						UserID:       1,
						AmountIDR:    150000,
						Description:  "Test Expense",
						Status:       int32(util.EXPENSE_PENDING),
						AutoApproved: false,
					}, nil).
					Times(1)

				server.MockRepo.EXPECT().
					ApprovalExpense(gomock.Any(), gomock.Any()).
					Return(nil).
					Times(1)

				server.MockPaymentProcessor.EXPECT().
					ProcessPayment(gomock.Any(), gomock.Any()).
					Return(&entity.PaymentProcessorResponse{
						Data: struct {
							ID         string `json:"id"`
							ExternalID string `json:"external_id"`
							Status     string `json:"status"`
						}{
							ID:         "TXN123",
							ExternalID: "EXT123",
							Status:     "SUCCESS",
						},
						Message: "Payment processed successfully",
					}, nil).
					Times(1)

				server.MockRepo.EXPECT().
					WriteAuditLog(gomock.Any(), gomock.Any()).
					Return(errors.New("audit log failed")).
					Times(1)
			},
			wantErr: false, // Should not fail even if audit log fails
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := NewTestServerWithPaymentProcessor(t)
			defer server.MockCtrl.Finish()

			ctx := context.Background()
			tt.mock(server)

			err := server.Service.ProcessPayment(ctx, tt.request)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			assert.NoError(t, err)
		})
	}
}
