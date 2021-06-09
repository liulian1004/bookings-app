package dbrepo

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/tsawler/bookings-app/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}
// insert a reservation into db
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3 seconds then cancel
	defer cancel()
	var newID int
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
	values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`

	err := m.DB.QueryRowContext(ctx,stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan((&newID))

	if err != nil {
		return 0,err
	}

	return newID, nil
}

//insert a room restriction into db
func (m *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3 seconds then cancel
	defer cancel()
	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, 
		   created_at, updated_at, restriction_id)
		   values($1,$2,$3,$4,$5,$6,$7)`
	_, err := m.DB.ExecContext(ctx, stmt,
			res.StartDate,
			res.EndDate,
			res.RoomID,
			res.ReservationID,
			time.Now(),
			time.Now(),
			res.RestrictionID,
	)

	if err != nil {
		return err
	}
	return nil
}
// search room availability for roomID
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start, end time.Time, roomID int) (bool, error) { 
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3 seconds then cancel
	defer cancel()
	//query for a certain room
	query := `select count(id) from room_restrictions where 
	room_id = $1 and
	$2 < end_date and $3 > start_date;`
	row := m.DB.QueryRowContext(ctx,query, roomID, start, end)
	var numRows int
	err := row.Scan(&numRows)
	if err != nil {
		log.Println("err:", err)
		return false, err
	}
	if numRows != 0 {
		return false, nil
	}

	return true, nil
}
// return the slice of avaiable rooms for given date
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3 seconds then cancel
	defer cancel()
	var rooms []models.Room
	query := `select r.id, r.room_name
			from rooms r
			where r.id not in
			(select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);`
	
	rows, err := m.DB.QueryContext(ctx,query, start, end)
	
	if err != nil {
		return rooms, err
	}

	for rows.Next(){
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err = rows.Err(); err != nil{
		return rooms, err
	}
	return rooms, nil
}

func (m *postgresDBRepo) GetRoomByID(id int) (models.Room, error) {
	//connect to db
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3 seconds then cancel
	defer cancel()

	var room models.Room

	query := `select id, room_name, created_at, updated_at from rooms where id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)

	if err != nil {
		return room, err
	}
	return room, nil
}

func (m *postgresDBRepo) GetuserByID(ID int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3 seconds then cancel
	defer cancel()

	query := `
		select id, first_name, last_name, email, password, access_level, created_at,
		updated_at from users where id = $1
	`
	row := m.DB.QueryRowContext(ctx, query, ID)

	var u models.User

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	return u, err 
}

func (m *postgresDBRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3 seconds then cancel
	defer cancel()

	query := `
		update users set first_name = $1, last_name = $2, email = $3, access_level = $4, update_at = $5
	`
	_, err := m.DB.ExecContext(ctx, query, 
		u.FirstName,
		u.LastName,
		u.Email,
		u.AccessLevel,
		time.Now(),
	)
	return err
}

func (m *postgresDBRepo) Authenticate(email, password string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second) // 3 seconds then cancel
	defer cancel()
	var id int
	var hashedPassWord string

	query := `select id, password from users where email = $1 `
	row := m.DB.QueryRowContext(ctx, query, email)

	err := row.Scan(&id, &hashedPassWord)

	if err != nil {
		return id, "", err
	}
	//compare the pw with the db ones
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassWord), []byte(password))
	

	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0,"", errors.New("incorrect password")
	}else if err != nil{
		return 0,"" ,err

	}
	return id, hashedPassWord, nil
}