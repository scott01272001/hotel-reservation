package types

type HotelQueryParams struct {
	Rooms bool
}

type Hotel struct {
	ID       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string   `bson:"name" json:"name"`
	Location string   `bson:"location" json:"location"`
	Rooms    []string `bson:"rooms" json:"rooms"`
	Rating   int      `bson:"rating" json:"rating"`
}

type RoomType int

const (
	SingleRoomType RoomType = iota + 1
	DoubleRoomType
	SeaSideRoomType
	DeluxRoomType
)

type Room struct {
	ID        string   `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType `bson:"type" json:"type"`
	BasePrice float64  `bson:"basePrice" json:"basePrice"`
	Price     float64  `bson:"Price" json:"Price"`
	HotelID   string   `bson:"hotelId" json:"hotelId"`
}
