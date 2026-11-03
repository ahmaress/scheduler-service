package models

import (
	"encoding/json"
	"time"
)

type AvailabilityRule struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Type           string    `json:"type,omitempty"`
	Event          string    `json:"event,omitempty"` // Event identifier to group availability rules
	DaysOfWeek     []int     `json:"days_of_week"`    // Array of days (0=Sunday, 6=Saturday), 1-7 days required
	StartTime      string    `json:"start_time"`
	EndTime        string    `json:"end_time"`
	SlotLengthMins int       `json:"slot_length_minutes"`
	Title          string    `json:"title,omitempty"`
	Available      bool      `json:"available"`
	IsRecurring    bool      `json:"is_recurring"` // true = all future weeks, false = current week only
	CreatedAt      time.Time `json:"created_at_utc,omitempty"`
	UpdatedAt      time.Time `json:"updated_at_utc,omitempty"`
}

// MarshalJSON ensures timestamps are serialized in UTC
func (a AvailabilityRule) MarshalJSON() ([]byte, error) {
	type Alias AvailabilityRule
	return json.Marshal(&struct {
		CreatedAtUTC time.Time `json:"created_at_utc,omitempty"`
		UpdatedAtUTC time.Time `json:"updated_at_utc,omitempty"`
		*Alias
	}{
		CreatedAtUTC: a.CreatedAt.UTC(),
		UpdatedAtUTC: a.UpdatedAt.UTC(),
		Alias:        (*Alias)(&a),
	})
}
