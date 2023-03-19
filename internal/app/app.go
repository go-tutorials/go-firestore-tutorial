package app

import (
	"context"
	"firebase.google.com/go"
	"github.com/core-go/health"
	"github.com/core-go/health/firestore"
	"google.golang.org/api/option"

	"go-service/internal/handler"
	"go-service/internal/service"
)

type ApplicationContext struct {
	Health *health.Handler
	User   *handler.UserHandler
}

func NewApp(ctx context.Context, cfg Config) (*ApplicationContext, error) {
	opts := option.WithCredentialsJSON([]byte(cfg.Credentials))
	app, err := firebase.NewApp(ctx, nil, opts)
	if err != nil {
		return nil, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, err
	}

	userService := service.NewUserService(client)
	userHandler := handler.NewUserHandler(userService)

	firestoreChecker := firestore.NewHealthChecker(ctx, []byte(cfg.Credentials))
	healthHandler := health.NewHandler(firestoreChecker)

	return &ApplicationContext{
		Health: healthHandler,
		User:   userHandler,
	}, nil
}
