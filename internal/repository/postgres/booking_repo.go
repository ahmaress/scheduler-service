package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"scheduler-service/internal/models"
	"scheduler-service/internal/repository"
)

type BookingRepo struct{}

func NewBookingRepo() *BookingRepo { return &BookingRepo{} }

func (r *BookingRepo) ListBookingsInRange(ctx context.Context, q repository.Querier, userID string, from, to repository.AppTime) ([]models.Booking, error) {
	query := `SELECT id,user_id,candidate_email,start_at_utc,end_at_utc,status,created_at 
		      FROM bookings
		      WHERE user_id=$1 AND start_at_utc >= $2 AND start_at_utc < $3 AND status='confirmed'`
	rows, err := q.Query(ctx, query, userID, from, to)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Booking
	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.CandidateEmail, &b.StartAtUTC, &b.EndAtUTC, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

func (r *BookingRepo) ListBookings(ctx context.Context, q repository.Querier, userID string, from, to repository.AppTime, filtered bool) ([]models.Booking, error) {
	var (
		rows pgx.Rows
		err  error
	)
	if filtered {
		query := `SELECT id,user_id,candidate_email,start_at_utc,end_at_utc,status,created_at 
		          FROM bookings 
		          WHERE user_id=$1 AND start_at_utc >= $2 AND start_at_utc < $3 AND status != 'cancelled'
		          ORDER BY start_at_utc`
		rows, err = q.Query(ctx, query, userID, from, to)
	} else {
		query := `SELECT id,user_id,candidate_email,start_at_utc,end_at_utc,status,created_at 
		          FROM bookings 
		          WHERE user_id=$1 AND status != 'cancelled'
		          ORDER BY start_at_utc`
		rows, err = q.Query(ctx, query, userID)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []models.Booking
	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(&b.ID, &b.UserID, &b.CandidateEmail, &b.StartAtUTC, &b.EndAtUTC, &b.Status, &b.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, b)
	}
	return out, nil
}

func (r *BookingRepo) CheckExistingBookingAtStart(ctx context.Context, q repository.Querier, userID string, start repository.AppTime) (string, error) {
	query := `SELECT id FROM bookings 
		       WHERE user_id=$1 AND status='confirmed' 
		       AND start_at_utc = $2 FOR UPDATE`
	var id string
	err := q.QueryRow(ctx, query, userID, start).Scan(&id)
	if err != nil && err != pgx.ErrNoRows {
		return "", err
	}
	return id, err
}
