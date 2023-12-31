package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/scott/hotel-reservation/types"
)

func HandlerGetUsers(c *fiber.Ctx) error {
	u := types.User{
		FirstName: "James",
		LastName:  "Bob",
	}
	return c.JSON(u)
}

func HandlerGetUser(c *fiber.Ctx) error {
	return c.JSON(fmt.Sprintf("id: %s, James", c.Params("id")))
}
