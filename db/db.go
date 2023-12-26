package db

const (
	DBURI      = "mongodb://localhost:27017/?tls=false"
	DBName     = "hotel-reservation"
	TESTDBName = "hotel-reservation-test"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
