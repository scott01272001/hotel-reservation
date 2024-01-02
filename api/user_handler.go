package api

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/scott/hotel-reservation/db"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlerGetUser(c *fiber.Ctx) error {
	var (
		id  = c.Params("id")
		ctx = context.Background()
	)
	user, err := h.userStore.GetUserById(ctx, id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return c.JSON(user)
}

func (h *UserHandler) HandlerGetUsers(c *fiber.Ctx) error {
	var (
		ctx = context.Background()
	)
	users := h.userStore.GetUsers(ctx)

	fmt.Println(users)

	return c.JSON(users)
}
