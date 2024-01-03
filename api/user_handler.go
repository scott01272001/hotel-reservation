package api

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/scott/hotel-reservation/db"
	"github.com/scott/hotel-reservation/types"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var user types.User
	err := c.BodyParser(&user)
	if err != nil {
		return err
	}
	h.userStore.PostUser(context.Background(), &user)
	return nil
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
	users, err := h.userStore.GetUsers(ctx)
	if err != nil {
		return err
	}

	fmt.Println(users)

	return c.JSON(users)
}
