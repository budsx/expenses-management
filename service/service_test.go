package service

import (
	"testing"

	repo "github.com/budsx/expenses-management/repository"
	_interface "github.com/budsx/expenses-management/repository/interface"
	"github.com/budsx/expenses-management/util"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
)

type TestService struct {
	MockCtrl             *gomock.Controller
	MockRepo             *_interface.MockExpensesRepository
	MockRabbitMQ         *_interface.MockRabbitMQClient
	MockUserRepo         *_interface.MockUserRepository
	MockPaymentProcessor *_interface.MockPaymentProcessor
	MockLogger           *logrus.Logger
	Service              *ExpensesManagementService
}

func NewTestServer(t *testing.T) *TestService {
	ctrl := gomock.NewController(t)
	mockRepo := _interface.NewMockExpensesRepository(ctrl)
	mockRabbitMQ := _interface.NewMockRabbitMQClient(ctrl)
	mockLogger := util.NewLogger(-1)
	service := NewExpensesManagementService(&repo.Repository{
		ExpensesRepository: mockRepo,
		RabbitMQClient:     mockRabbitMQ,
	}, mockLogger)

	return &TestService{
		MockCtrl:     ctrl,
		MockRepo:     mockRepo,
		MockRabbitMQ: mockRabbitMQ,
		MockLogger:   mockLogger,
		Service:      service,
	}
}

func NewTestServerWithUserRepo(t *testing.T) *TestService {
	ctrl := gomock.NewController(t)
	mockRepo := _interface.NewMockExpensesRepository(ctrl)
	mockRabbitMQ := _interface.NewMockRabbitMQClient(ctrl)
	mockUserRepo := _interface.NewMockUserRepository(ctrl)
	mockLogger := util.NewLogger(-1)
	service := NewExpensesManagementService(&repo.Repository{
		ExpensesRepository: mockRepo,
		RabbitMQClient:     mockRabbitMQ,
		UserRepository:     mockUserRepo,
	}, mockLogger)

	return &TestService{
		MockCtrl:     ctrl,
		MockRepo:     mockRepo,
		MockRabbitMQ: mockRabbitMQ,
		MockUserRepo: mockUserRepo,
		MockLogger:   mockLogger,
		Service:      service,
	}
}

func NewTestServerWithPaymentProcessor(t *testing.T) *TestService {
	ctrl := gomock.NewController(t)
	mockRepo := _interface.NewMockExpensesRepository(ctrl)
	mockRabbitMQ := _interface.NewMockRabbitMQClient(ctrl)
	mockUserRepo := _interface.NewMockUserRepository(ctrl)
	mockPaymentProcessor := _interface.NewMockPaymentProcessor(ctrl)
	mockLogger := util.NewLogger(-1)
	service := NewExpensesManagementService(&repo.Repository{
		ExpensesRepository: mockRepo,
		RabbitMQClient:     mockRabbitMQ,
		UserRepository:     mockUserRepo,
		PaymentProcessor:   mockPaymentProcessor,
	}, mockLogger)

	return &TestService{
		MockCtrl:             ctrl,
		MockRepo:             mockRepo,
		MockRabbitMQ:         mockRabbitMQ,
		MockUserRepo:         mockUserRepo,
		MockPaymentProcessor: mockPaymentProcessor,
		MockLogger:           mockLogger,
		Service:              service,
	}
}
