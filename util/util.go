package util

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/budsx/expenses-management/model"
)

func GetUserInfoFromContext(ctx context.Context) (model.User, error) {
	userID, ok := ctx.Value("user_id").(int64)
	if !ok {
		return model.User{}, fmt.Errorf("user_id not found")
	}

	email, ok := ctx.Value("user_email").(string)
	if !ok {
		return model.User{}, fmt.Errorf("email not found")
	}

	role, ok := ctx.Value("user_role").(int)
	if !ok {
		return model.User{}, fmt.Errorf("role not found")
	}

	return model.User{
		ID:    userID,
		Email: email,
		Role:  role,
	}, nil
}

func OnShutdown(shutdown func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)
	<-c
	id := time.Now().UnixNano()
	fmt.Println("OnShutdown...", id)
	if shutdown != nil {
		shutdown()
	}
	fmt.Println("OnShutdown done", id)
}
