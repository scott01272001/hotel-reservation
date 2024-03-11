package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/scott/hotel-reservation/db"
	"github.com/scott/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserHandler struct {
	store *db.Store
}

func NewUserHandler(store *db.Store) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *UserHandler) HandleAuthenticate(c *fiber.Ctx) error {
	var authParams AuthParams
	if err := c.BodyParser(&authParams); err != nil {
		return err
	}

	fmt.Println(authParams)

	return nil
}

func (h *UserHandler) HandleDeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.store.User.DeleteUser(c.Context(), id); err != nil {
		return err
	}
	return c.JSON(map[string]string{"deleted": id})
}

func (h *UserHandler) HandlePutUser(c *fiber.Ctx) error {
	var (
		params types.UpdateUserParams
		userId = c.Params("id")
	)
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	errors, err := types.ValidateStruct(&params)
	if err != nil {
		return err
	}
	if len(errors) != 0 {
		return c.JSON(errors)
	}

	updated, err := h.store.User.UpdateUser(c.Context(), userId, &params)
	if err != nil {
		return err
	}
	return c.JSON(updated)
}

func (h *UserHandler) HandlePostUser(c *fiber.Ctx) error {
	var params types.CreateUserParam
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	errors, err := types.ValidateStruct(&params)
	if err != nil {
		return err
	}
	if len(errors) != 0 {
		return c.JSON(errors)
	}

	user, err := types.NewUserFromParams(&params)
	if err != nil {
		return err
	}

	insertedUser, err := h.store.User.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandlerGetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.store.User.GetUserById(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
		return err
	}
	return c.JSON(user)
}

func (h *UserHandler) HandlerGetUsers(c *fiber.Ctx) error {
	users, err := h.store.User.GetUsers(c.Context())
	if err != nil {
		return err
	}
	return c.JSON(users)
}
