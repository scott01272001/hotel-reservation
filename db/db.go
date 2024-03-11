package db

const DBNAME = "hotel-reservation"

const DBURI = "mongodb+srv://Cluster28936:YkhTcXV+bHBl@hotel-reservation.jtvk28l.mongodb.net/?retryWrites=true&w=majority"

type Store struct {
	User  UserStore
	Hotel HotelStore
	Room  RoomStore
}
