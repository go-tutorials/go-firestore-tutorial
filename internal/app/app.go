package app

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/core-go/health"
	"github.com/core-go/health/firestore"
	"google.golang.org/api/option"

	"go-service/internal/handlers"
	"go-service/internal/services"
)

type ApplicationContext struct {
	HealthHandler *health.HealthHandler
	UserHandler   *handlers.UserHandler
}

func NewApp(ctx context.Context, root Root) (*ApplicationContext, error) {
	sa := option.WithCredentialsFile(root.Firestore.File)
	app, er1 := firebase.NewApp(ctx, nil, sa)
	if er1 != nil {
		return nil, er1
	}

	client, er2 := app.Firestore(ctx)
	if er2 != nil {
		return nil, er2
	}

	userService := services.NewUserService(client)
	userHandler := handlers.NewUserHandler(userService)

	firestoreChecker := firestore.NewHealthChecker(root.Firestore.ProjectId)
	healthHandler := health.NewHealthHandler(firestoreChecker)

	return &ApplicationContext{
		HealthHandler: healthHandler,
		UserHandler:   userHandler,
	}, nil
}
