package service

import (
	"time"

	"scheduler-service/internal/repository"
)

type AvailabilityService struct {
	DB    repository.Querier
	Avail repository.AvailabilityRepository
	Book  repository.BookingRepository
}

type Slot struct {
	StartUTC time.Time `json:"start_utc"`
	EndUTC   time.Time `json:"end_utc"`
	Event    string    `json:"event,omitempty"` // Event identifier for the slot
}

func NewAvailabilityService(db repository.Querier, ar repository.AvailabilityRepository, br repository.BookingRepository) *AvailabilityService {
	return &AvailabilityService{DB: db, Avail: ar, Book: br}
}
