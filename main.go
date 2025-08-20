package main

import (
	"fmt"

	"github.com/budsx/expenses-management/config"
	"github.com/budsx/expenses-management/handler"
	"github.com/budsx/expenses-management/repository"
	"github.com/budsx/expenses-management/repository/paymentprocessor"
	"github.com/budsx/expenses-management/repository/postgres"
	"github.com/budsx/expenses-management/service"
	"github.com/budsx/expenses-management/transport/http"
	"github.com/budsx/expenses-management/util"
)

func main() {
	conf := config.Load()
	logger := util.NewLogger(conf.Log.Level)

	conn, err := config.NewDatabase(conf, logger)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize database")
		return
	}
	defer conn.Close()

	repos := repository.Repository{
		PaymentProcessor:   paymentprocessor.NewPaymentAPI("https://api.paymentprocessor.com"),
		UserRepository:     postgres.NewUserRepository(conn),
		ExpensesRepository: postgres.NewExpensesRepository(conn),
	}

	service := service.NewExpensesManagementService(&repos, logger)

	expensesHandler := handler.NewExpensesManagementHandler(service)
	userHandler := handler.NewUserHandler(service)

	server := http.NewExpensesManagementServer(service, userHandler, expensesHandler)
	server.Run(fmt.Sprintf(":%d", conf.ServicePort))
}
