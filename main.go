package main

import (
	"context"
	"flag"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/scott/hotel-reservation/api"
	"github.com/scott/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb+srv://Cluster28936:YkhTcXV+bHBl@hotel-reservation.jtvk28l.mongodb.net/?retryWrites=true&w=majority"
const dbname = "hotel-reservation"
const userColl = "users"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	listenAddr := flag.String("listenAddr", ":8081", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	api.NewUserHandler(db.NewMongoUserStore(client))

	app := fiber.New()
	apiv1 := app.Group("/api/v1")

	apiv1.Get("/user", api.HandlerGetUsers)
	apiv1.Get("/user/:id", api.HandlerGetUser)

	err = app.Listen(*listenAddr)
	if err != nil {
		log.Fatal(err)
	}
}
