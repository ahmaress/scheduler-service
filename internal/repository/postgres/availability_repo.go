package postgres

import (
	"context"
	"time"

	"scheduler-service/internal/models"
	"scheduler-service/internal/repository"
)

// normalizeTime converts database time format (HH:MM:SS.microseconds) to HH:MM format
func normalizeTime(timeStr string) string {
	// Handle empty string
	if timeStr == "" {
		return ""
	}
	// Take only first 5 characters (HH:MM)
	if len(timeStr) >= 5 {
		return timeStr[:5]
	}
	return timeStr
}

type AvailabilityRepo struct{}

func NewAvailabilityRepo() *AvailabilityRepo { return &AvailabilityRepo{} }

func (r *AvailabilityRepo) InsertAvailabilityRule(ctx context.Context, q repository.Querier, ar *models.AvailabilityRule) error {
	now := time.Now().UTC()
	query := `INSERT INTO availability_rules
		(id, user_id, type, event, days_of_week, start_time, end_time, slot_length_minutes, title, available, is_recurring, created_at, updated_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`
	var event *string
	if ar.Event != "" {
		event = &ar.Event
	}
	return q.QueryRow(ctx, query,
		ar.UserID, ar.Type, event, ar.DaysOfWeek, ar.StartTime, ar.EndTime, ar.SlotLengthMins,
		ar.Title, ar.Available, ar.IsRecurring, now, now,
	).Scan(&ar.ID)
}

