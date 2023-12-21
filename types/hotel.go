package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string               `bson:"name" json:"name"`
	Location string               `bson:"location" json:"location"`
	Rooms    []primitive.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int                  `bson:"rating" json:"rating"`
}

type RoomType int

const (
	SingleRoomType RoomType = iota + 1
	DoubleRoomType
	SeaSideRoomType
	DeluxeRoomType
)

func (t RoomType) String() string {
	//todo finish using swich
	return "SingleRoomType"
}

type Room struct {
	ID      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	HotelID primitive.ObjectID `bson:"hotel_id" json:"hotel_id"`
	//small,normal,kingsize
	Size      string  `bson:"size" json:"size"`
	SeaSide   bool    `bson:"sea_side" json:"sea_side"`
	Price     float64 `bson:"price" json:"price"`
}
