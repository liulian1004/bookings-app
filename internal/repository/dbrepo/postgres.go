package dbrepo

import (
	"context"
	"time"

	"github.com/tsawler/bookings-app/internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}
// insert a reservation into db
func (m *postgresDBRepo) InsertReservation(res models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3 seconds then cancel
	defer cancel()
	
	stmt := `insert into reservations (first_name, last_name, email, phone, start_data, end_data, room_id, created_at, updated_at)
	value($1,$2,$3,$4,$5,$6,$7,$8, $9)`

	_, err := m.DB.ExecContext(ctx,stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}