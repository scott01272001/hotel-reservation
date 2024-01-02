package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/scott/hotel-reservation/api"
	"github.com/scott/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb+srv://Cluster28936:YkhTcXV+bHBl@hotel-reservation.jtvk28l.mongodb.net/?retryWrites=true&w=majority"

func errorHandler(ctx *fiber.Ctx, err error) error {
	return ctx.JSON(map[string]string{"error": err.Error()})
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	listenAddr := flag.String("listenAddr", ":8082", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	userHandler := api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New(fiber.Config{
		ErrorHandler: errorHandler,
	})
	app.Use(recover.New())

	apiv1 := app.Group("/api/v1")
	apiv1.Get("/users", userHandler.HandlerGetUsers)
	apiv1.Get("/users/:id", userHandler.HandlerGetUser)

	apiv1.Get("/error", func(c *fiber.Ctx) error {
		panic("boom!!!!")
		return nil
	})

	err = app.Listen(*listenAddr)
	if err != nil {
		log.Fatal(err)
	}
}
