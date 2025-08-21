package service

import (
	"context"
	"errors"
	"testing"

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
