package main

import (
	"context"
	"hexagonal-gotest/handlers"
	"hexagonal-gotest/repositories"
	"hexagonal-gotest/services"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initMongo() *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://ep5-course:HlT9NpyD4Vt0HtbK@cluster0.vvx397a.mongodb.net/?retryWrites=true&w=majority"))
	if err != nil {
		panic(err)
	}
	// # Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}
	return client.Database("julladith")
}

func main() {
	db := initMongo()

	//init Data Layer
	userRepo := repositories.NewUserRepository(db, "users")

	//init Business Logic Layer
	userSrv := services.NewUserService(userRepo)

	//init Presentation Layer
	userHand := handlers.NewUserHandler(userSrv)

	//framework routes
	app := fiber.New()
	app.Post("/register", userHand.Register)

	//start server
	app.Listen(":3000")
}
