package main

import (
	"context"
	"flag"
	"github.com/scott/hotel-reservation/api/middeware"
	"log"

	"github.com/gofiber/fiber/v2"
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

	listenAddr := flag.String("listenAddr", ":8032", "The listen address of the API server")
	flag.Parse()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
	}

	var (
		// store init
		hotelStore = db.NewMongoHotelStore(client)
		roomStore  = db.NewMongoRoomStore(client, hotelStore)
		userStore  = db.NewMongoUserStore(client)
		store      = db.Store{
			Hotel: hotelStore,
			Room:  roomStore,
			User:  userStore,
		}

		// handlers init
		userHandler  = api.NewUserHandler(&store)
		hotelHandler = api.NewHotelHandler(&store)
		app          = fiber.New(fiber.Config{
			ErrorHandler: errorHandler,
		})

		apiNoAuth = app.Group("/api/")
		apiv1     = app.Group("/api/v1", middeware.JWTAuthentication)
	)

	// auth handlers
	apiNoAuth.Post("/auth", userHandler.HandleAuthenticate)

	// user handlers
	apiv1.Get("/users", userHandler.HandlerGetUsers)
	apiv1.Get("/users/:id", userHandler.HandlerGetUser)
	apiv1.Post("/users", userHandler.HandlePostUser)
	apiv1.Delete("/users/:id", userHandler.HandleDeleteUser)
	apiv1.Put("/users/:id", userHandler.HandlePutUser)

	// hotel handlers
	apiv1.Get("/hotels", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotels/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotels/:id/rooms", hotelHandler.HandleGetRooms)

	err = app.Listen(*listenAddr)
	if err != nil {
		log.Fatal(err)
	}
}
