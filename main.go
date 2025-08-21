package main

import (
	"fmt"

	"github.com/budsx/expenses-management/config"
	"github.com/budsx/expenses-management/handler"
	"github.com/budsx/expenses-management/repository"
	"github.com/budsx/expenses-management/repository/payment"
	"github.com/budsx/expenses-management/repository/postgres"
	"github.com/budsx/expenses-management/repository/rabbitmq"
	"github.com/budsx/expenses-management/service"
	"github.com/budsx/expenses-management/transport/http"
	"github.com/budsx/expenses-management/transport/messaging"
	"github.com/budsx/expenses-management/util"
	utilRabbitmq "github.com/budsx/expenses-management/util/rabbitmq"
)

func main() {
	conf := config.Load()
	logger := util.NewLogger(conf.Log.Level)
	conn, err := config.NewDatabase(conf, logger)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize database")
		return
	}

	rabbitmqClient, err := utilRabbitmq.NewClient(conf.RabbitMQURL)
	if err != nil {
		logger.WithError(err).Error("Failed to initialize rabbitmq client")
		return
	}

	repos := repository.NewRepository(
		payment.NewPaymentProcessor(conf.PaymentProcessorURL),
		postgres.NewUserRepository(conn),
		postgres.NewExpensesRepository(conn),
		rabbitmq.NewRabbitClient(rabbitmqClient, conf.TopicPaymentProcessor),
	)
	service := service.NewExpensesManagementService(repos, logger)
	expensesHandler := handler.NewExpensesManagementHandler(service)
	authHandler := handler.NewAuthHandler(service)

	server := http.NewExpensesManagementServer(service, expensesHandler, authHandler)

	messaging.NewTransportListener(service, repos.RabbitMQClient, conf.TopicPaymentProcessor, conf.TopicPaymentProcessor+".ems.queue")
	server.ServeHTTP(fmt.Sprintf(":%d", conf.ServicePort))
	logger.Info("Server started...")

	util.OnShutdown(func() {
		logger.Info("Shutting down server...")
		server.Shutdown()
		logger.Info("Server shutdown...")
		conn.Close()
		logger.Info("connection closed...")
	})
}
