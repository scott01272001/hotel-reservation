package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/scott/hotel-reservation/db"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore *db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: *userStore,
	}
}

func (h *UserHandler) HandlerGetUsers(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userStore.GetUserById(id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandlerGetUser(c *fiber.Ctx) error {
	return c.JSON(fmt.Sprintf("id: %s, James", c.Params("id")))
}
