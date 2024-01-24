package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/scott/hotel-reservation/db"
	"github.com/scott/hotel-reservation/types"
)

type HotelHandler struct {
	roomStore  db.RoomStore
	hotelStore db.HotelStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler {
	return &HotelHandler{
		hotelStore: hs,
		roomStore:  rs,
	}
}

func (h *HotelHandler) HandleGetHotels(c *fiber.Ctx) error {
	var qparams types.HotelQueryParams
	if err := c.QueryParser(&qparams); err != nil {
		return err
	}
	fmt.Println(qparams)

	hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
	if err != nil {
		return err
	}
	return c.JSON(hotels)
}
