package db

var (
	MONGODBENVDBNAME = "MONGO_DB_NAME"
	MONGODBENVDBURI  = "MONGO_DB_URI"

	MONGOTESTDBENVDBNAME = "MONGO_TEST_DB_NAME"
	MONGOTESTDBENVDBURI  = "MONGO_TEST_DB_URI"
)

type Store struct {
	User    UserStore
	Hotel   HotelStore
	Room    RoomStore
	Booking BookingStore
}
