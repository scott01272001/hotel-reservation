package api

import (
	"context"
	"log"

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
	var params types.CreateUserparam
	if err := c.BodyParser(&params); err != nil {
		log.Fatal(err)
		return err
	}

	if err := types.ValidateStruct(params); err != nil {
		return c.JSON(err.Error())
	}

	user, err := types.NewUserFromParams(&params)
	if err != nil {
		log.Fatal(err)
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return c.JSON(insertedUser)
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
	return c.JSON(users)
}
