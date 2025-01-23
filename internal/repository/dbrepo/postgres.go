package dbrepo

import (
	"context"
	"time"

	"github.com/janomonje/bed-n-breakfast/internal/models"
)

func (model *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation inserts a reservation into the databse
func (model *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	statement := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at) values ($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	err := model.DB.QueryRowContext(ctx, statement,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now()).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

// InsertRoomRestriction inserts a room restriction into the database
func (model *postgresDBRepo) InsertRoomRestriction(res models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	statement := `insert into room_restrictions(start_date, end_date, room_id, reservation_id,
			      created_at, updated_at, restriction_id)
				  values
				  ($1,$2,$3,$4,$5,$6,$7)`
	_, err := model.DB.ExecContext(ctx, statement,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
		time.Now(),
		time.Now(),
		res.RestrictionID)

	if err != nil {
		return err
	}
	return nil

}

// SearchAvailabilityByDates returns true if availability exists for roomID, and false if no availability exists
func (model *postgresDBRepo) SearchAvailabilityByDates(start, end time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var numRows int
	query := `
	select 
		count(id)
	from
		room_restrictions 
	where 
		room_id = $1
		$2 < end_date and $3 > start_date`

	row := model.DB.QueryRowContext(ctx, query, roomID,start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}
