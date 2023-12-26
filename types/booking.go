package types

import (
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookingParam struct {
	FromDate   time.Time `json:"from_date"`
	TillDate   time.Time `json:"till_date"`
	NumPersons int       `json:"num_persons"`
}

func (b BookingParam) Validate() error {
	errs := []error{}
	now := time.Now()

	if b.FromDate.Before(now) || b.TillDate.Before(now) {
		errs = append(errs, fmt.Errorf("cannot book in the past"))
	}

	return errors.Join(errs...)
}

type Booking struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	RoomID     primitive.ObjectID `bson:"room_id" json:"room_id"`
	FromDate   time.Time          `bson:"from_date" json:"from_date"`
	TillDate   time.Time          `bson:"till_date" json:"till_date"`
	NumPersons int                `bson:"num_persons" json:"num_persons"`
}
