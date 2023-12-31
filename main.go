package main

import (
	"flag"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/scott/hotel-reservation/api"
)

func main() {
	listenAddr := flag.String("listenAddr", ":8081", "The listen address of the API server")
	flag.Parse()

	app := fiber.New()

	apiv1 := app.Group("/api/v1")
	apiv1.Get("/user", api.HandlerGetUsers)
	apiv1.Get("/user/:id", api.HandlerGetUser)

	err := app.Listen(*listenAddr)
	if err != nil {
		fmt.Println(err.Error())
	}
}
