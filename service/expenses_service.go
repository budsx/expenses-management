package service

import (
	"context"
	"fmt"
	"time"

	"github.com/budsx/expenses-management/entity"
	"github.com/budsx/expenses-management/model"
	"github.com/budsx/expenses-management/util"
	"github.com/google/uuid"
)

func (s *ExpensesManagementService) CreateExpense(ctx context.Context, req model.CreateExpenseRequest) (*model.ExpenseResponse, error) {
	userInfo, err := util.GetUserInfoFromContext(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user info")
		return nil, fmt.Errorf("failed to get user info")
	}

	s.logger.WithField("user_id", userInfo.ID).Info("User is submitting expense")

	var autoApproved bool
	if req.AmountIDR < 1000000 {
		autoApproved = true
	}

	expenseID, err := s.repo.ExpensesRepository.WriteExpense(ctx, &entity.Expense{
		UserID:      userInfo.ID,
		AmountIDR:   req.AmountIDR,
		Description: req.Description,
		ReceiptURL:  req.ReceiptURL,
		Status:      int32(util.EXPENSE_PENDING),
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to write expense")
		return nil, err
	}

	if autoApproved {
		util.GoWithRecover(func() {
			err = s.repo.RabbitMQClient.PublishPayment(&entity.PublishPaymentRequest{
				ExpenseID:  expenseID,
				ApproverID: 0, // Auto approved, no approver
				Notes:      "Auto Approved",
				Status:     int32(util.EXPENSE_AUTO_APPROVED),
			})
			if err != nil {
				s.logger.WithError(err).Error("failed to publish payment")
			}
			s.logger.Info("Publish to payment processor")
		})
	}

	err = s.repo.ExpensesRepository.WriteAuditLog(ctx, &entity.AuditLog{
		ExpenseID:    expenseID,
		NewStatus:    int32(util.EXPENSE_PENDING),
		StatusBefore: int32(util.EXPENSE_PENDING),
		Notes:        "Expense created",
		CreatedAt:    time.Now(),
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to write audit log")
	}

	return &model.ExpenseResponse{
		ID:           expenseID,
		UserID:       userInfo.ID,
		AmountIDR:    req.AmountIDR,
		Description:  req.Description,
		ReceiptURL:   req.ReceiptURL,
		Status:       util.GetExpenseStatusString(util.EXPENSE_PENDING),
		AutoApproved: autoApproved,
	}, nil
}

func (s *ExpensesManagementService) GetExpenses(ctx context.Context, query model.ExpenseListQuery) (*model.ExpenseListResponse, error) {
	s.logger.WithField("query", query).Info("GetExpenses")
	userInfo, err := util.GetUserInfoFromContext(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user info")
		return nil, fmt.Errorf("failed to get user info")
	}

	if userInfo.Role == int(util.USER_ROLE_EMPLOYEE) {
		query.UserID = userInfo.ID
	}

	expenses, total, err := s.repo.ExpensesRepository.GetExpensesWithPagination(ctx, &entity.ExpenseListQuery{
		Page:   int32(query.Page),
		Limit:  int32(query.PageSize),
		UserID: query.UserID,
		Status: int32(query.Status),
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to get expenses")
		return nil, err
	}

	expensesResponse := make([]model.ExpenseResponse, 0)
	for _, expense := range expenses {
		expensesResponse = append(expensesResponse, model.ExpenseResponse{
			ID:           expense.ID,
			UserID:       expense.UserID,
			AmountIDR:    expense.AmountIDR,
			Description:  expense.Description,
			ReceiptURL:   expense.ReceiptURL,
			Status:       util.GetExpenseStatusString(util.ExpenseStatus(expense.Status)),
			AutoApproved: expense.AutoApproved,
		})
	}

	return &model.ExpenseListResponse{
		Expenses: expensesResponse,
		Total:    total,
		Page:     query.Page,
		PageSize: query.PageSize,
	}, nil
}

func (s *ExpensesManagementService) GetExpenseByID(ctx context.Context, expenseID int64) (*model.ExpenseResponse, error) {
	s.logger.WithField("expense_id", expenseID).Info("GetExpenseByID")
	expense, err := s.repo.ExpensesRepository.GetExpenseByID(ctx, expenseID)
	if err != nil {
		s.logger.WithError(err).Error("failed to get expense")
		return nil, err
	}

	return &model.ExpenseResponse{
		ID:           expense.ID,
		UserID:       expense.UserID,
		AmountIDR:    expense.AmountIDR,
		Description:  expense.Description,
		ReceiptURL:   expense.ReceiptURL,
		Status:       util.GetExpenseStatusString(util.ExpenseStatus(expense.Status)),
		AutoApproved: expense.AutoApproved,
	}, nil
}

func (s *ExpensesManagementService) ApproveExpense(ctx context.Context, req model.ApprovalRequest) (*model.ApprovalResponse, error) {
	userInfo, err := util.GetUserInfoFromContext(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user info")
		return nil, fmt.Errorf("failed to get user info")
	}

	if userInfo.Role != int(util.USER_ROLE_MANAGER) {
		s.logger.WithField("user_id", userInfo.ID).Error("user is not a manager")
		return nil, fmt.Errorf("user is not a manager")
	}

	util.GoWithRecover(func() {
		err = s.repo.RabbitMQClient.PublishPayment(&entity.PublishPaymentRequest{
			ExpenseID:  req.ExpenseID,
			ApproverID: userInfo.ID,
			Notes:      req.Notes,
			Status:     int32(util.APPROVAL_APPROVED),
		})
		if err != nil {
			s.logger.WithError(err).Error("failed to publish payment")
		}
	})

	return &model.ApprovalResponse{
		Message: fmt.Sprintf("Expense %d approved", req.ExpenseID),
	}, nil
}

func (s *ExpensesManagementService) RejectExpense(ctx context.Context, req model.ApprovalRequest) (*model.ApprovalResponse, error) {
	userInfo, err := util.GetUserInfoFromContext(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user info")
		return nil, fmt.Errorf("failed to get user info")
	}

	if userInfo.Role != int(util.USER_ROLE_MANAGER) {
		s.logger.WithField("user_id", userInfo.ID).Error("user is not a manager")
		return nil, fmt.Errorf("user is not a manager")
	}

	err = s.repo.ExpensesRepository.ApprovalExpense(ctx, &entity.ExpenseApproval{
		ExpenseID:  req.ExpenseID,
		ApproverID: userInfo.ID,
		Status:     int32(util.EXPENSE_REJECTED),
		Notes:      req.Notes,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to reject expense")
		return nil, fmt.Errorf("failed to reject expense")
	}

	err = s.repo.ExpensesRepository.WriteAuditLog(ctx, &entity.AuditLog{
		ExpenseID:    req.ExpenseID,
		NewStatus:    int32(util.EXPENSE_REJECTED),
		StatusBefore: int32(util.EXPENSE_PENDING),
		Notes:        req.Notes,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to write audit log")
	}

	return &model.ApprovalResponse{
		Message: fmt.Sprintf("Expense %d successfully rejected", req.ExpenseID),
	}, nil
}

func (s *ExpensesManagementService) HealthCheck(ctx context.Context) error {
	err := s.repo.ExpensesRepository.PingContext(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to ping database")
		return err
	}
	return nil
}

func (s *ExpensesManagementService) ProcessPayment(ctx context.Context, req model.ApprovalRequest) error {
	s.logger.WithField("request", req).Info("ProcessPayment")

	expense, err := s.repo.ExpensesRepository.GetExpenseByID(ctx, req.ExpenseID)
	if err != nil {
		s.logger.WithError(err).Error("failed to get expense")
		return fmt.Errorf("failed to get expense")
	}

	if expense.Status != int32(util.EXPENSE_PENDING) && expense.Status != int32(util.EXPENSE_AUTO_APPROVED) {
		s.logger.WithField("expense_id", req.ExpenseID).Error("expense is not pending")
		return fmt.Errorf("expense is not pending")
	}

	err = s.repo.ExpensesRepository.ApprovalExpense(ctx, &entity.ExpenseApproval{
		ExpenseID:  req.ExpenseID,
		ApproverID: req.ApproverID,
		Status:     req.Status,
		Notes:      req.Notes,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to approve expense")
		return fmt.Errorf("failed to approve expense")
	}
	s.logger.WithField("expense_id", req.ExpenseID).Info("Expense approved")

	payment, err := s.repo.PaymentProcessor.ProcessPayment(ctx, &entity.PaymentProcessorRequest{
		AmountIDR:  int64(expense.AmountIDR),
		ExternalID: uuid.New().String(),
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to process payment")
		return fmt.Errorf("failed to process payment")
	}
	s.logger.WithField("response", payment).Info("Payment processed")

	err = s.repo.ExpensesRepository.WriteAuditLog(ctx, &entity.AuditLog{
		ExpenseID:    req.ExpenseID,
		NewStatus:    int32(req.Status),
		StatusBefore: expense.Status,
		Notes:        req.Notes,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to write audit log")
	}

	s.logger.WithField("expense_id", req.ExpenseID).Info("Expense processed")
	return nil
}
