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

	// TODO: Get Rule from config
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
	s.logger.WithField("query", query).Info("Getting expenses")
	userInfo, err := util.GetUserInfoFromContext(ctx)
	if err != nil {
		s.logger.WithError(err).Error("failed to get user info")
		return nil, fmt.Errorf("failed to get user info")
	}

	if userInfo.Role != int(util.USER_ROLE_MANAGER) {
		s.logger.WithField("user_id", userInfo.ID).Error("user is not a manager")
		return nil, fmt.Errorf("user is not a manager")
	}

	expenses, total, err := s.repo.ExpensesRepository.GetExpensesWithPagination(ctx, &entity.ExpenseListQuery{
		Page:   int32(query.Page),
		Limit:  int32(query.PageSize),
		UserID: query.UserID,
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
	s.logger.WithField("expense_id", expenseID).Info("Getting expense by ID")
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

	externalID := uuid.New().String()
	payment, err := s.repo.PaymentProcessor.ProcessPayment(ctx, &entity.PaymentProcessorRequest{
		AmountIDR:  req.ExpenseID,
		ExternalID: externalID,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to process payment")
		return nil, fmt.Errorf("failed to process payment")
	}

	s.logger.WithField("Payment Processor Response", payment).Info("Payment processed")

	err = s.repo.ExpensesRepository.ApprovalExpense(ctx, &entity.ExpenseApproval{
		ExpenseID:  req.ExpenseID,
		ApproverID: userInfo.ID,
		Status:     int32(util.EXPENSE_APPROVED),
		Notes:      req.Notes,
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to approve expense")
		return nil, fmt.Errorf("failed to approve expense")
	}

	err = s.repo.ExpensesRepository.WriteAuditLog(ctx, &entity.AuditLog{
		ExpenseID:    req.ExpenseID,
		NewStatus:    int32(util.EXPENSE_APPROVED),
		StatusBefore: int32(util.EXPENSE_PENDING),
		Notes:        req.Notes,
		CreatedAt:    time.Now(),
	})
	if err != nil {
		s.logger.WithError(err).Error("failed to write audit log")
	}

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
