package repository

import (
	"time"

	"github.com/tsawler/bookings-app/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	SearchAvailabilityByDatesByRoomID(start, end time.Time,roomID int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetuserByID(ID int) (models.User, error)
	UpdateUser(u models.User) error
	Authenticate(email, password string) (int, string, error)
	AllReservation() ([] models.Reservation, error)
	AllNewReservation() ([] models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error) {

}