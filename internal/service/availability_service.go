package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"scheduler-service/internal/models"
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

func (s *AvailabilityService) SetAvailability(ctx context.Context, userID string, rules []models.AvailabilityRule) ([]models.AvailabilityRule, error) {
	var saved []models.AvailabilityRule
	for i := range rules {
		rules[i].UserID = userID
		now := time.Now().UTC()
		rules[i].CreatedAt = now
		rules[i].UpdatedAt = now
		if err := validateAvailabilityRule(&rules[i]); err != nil {
			return nil, err
		}
		if err := s.Avail.InsertAvailabilityRule(ctx, s.DB, &rules[i]); err != nil {
			return nil, err
		}
		saved = append(saved, rules[i])
	}
	return saved, nil
}

func (s *AvailabilityService) UpdateAvailability(ctx context.Context, userID, ruleID string, rule *models.AvailabilityRule) (*models.AvailabilityRule, error) {
	// Fetch existing rule first
	existing, err := s.Avail.GetAvailabilityRule(ctx, s.DB, userID, ruleID)
	if err != nil {
		return nil, err
	}
	// If days_of_week is empty, preserve the original value
	if len(rule.DaysOfWeek) == 0 {
		rule.DaysOfWeek = existing.DaysOfWeek
	}
	// If start_time is empty, preserve the original value
	if rule.StartTime == "" {
		rule.StartTime = existing.StartTime
	}
	// If end_time is empty, preserve the original value
	if rule.EndTime == "" {
		rule.EndTime = existing.EndTime
	}
	// If slot_length_minutes is zero, preserve the original value
	if rule.SlotLengthMins == 0 {
		rule.SlotLengthMins = existing.SlotLengthMins
	}
	if err := validateAvailabilityRule(rule); err != nil {
		return nil, err
	}
	id, err := s.Avail.UpdateAvailabilityRule(ctx, s.DB, userID, ruleID, rule)
	if err != nil {
		return nil, err
	}
	// Fetch the updated record from database to get correct timestamps
	updatedRule, err := s.Avail.GetAvailabilityRule(ctx, s.DB, userID, id)
	if err != nil {
		return nil, err
	}
	return updatedRule, nil
}

func (s *AvailabilityService) DeleteAvailability(ctx context.Context, userID, ruleID string) error {
	return s.Avail.DeleteAvailabilityRule(ctx, s.DB, userID, ruleID)
}

func (s *AvailabilityService) ListAvailability(ctx context.Context, userID string) ([]models.AvailabilityRule, error) {
	return s.Avail.ListAvailabilityRules(ctx, s.DB, userID)
}

func (s *AvailabilityService) ListBookings(ctx context.Context, userID string, from, to time.Time, filtered bool) ([]models.Booking, error) {
	return s.Book.ListBookings(ctx, s.DB, userID, from, to, filtered)
}

func validateAvailabilityRule(rule *models.AvailabilityRule) error {
	// Validate days_of_week: must have between 1 and 7 days
	if len(rule.DaysOfWeek) < 1 {
		return errors.New("days_of_week must contain at least 1 day")
	}
	if len(rule.DaysOfWeek) > 7 {
		return errors.New("days_of_week must contain at most 7 days")
	}
	
	// Validate each day is between 0-6 (Sunday=0 to Saturday=6)
	seenDays := make(map[int]bool)
	for _, day := range rule.DaysOfWeek {
		if day < 0 || day > 6 {
			return fmt.Errorf("invalid day_of_week: %d (must be between 0-6)", day)
		}
		if seenDays[day] {
			return fmt.Errorf("duplicate day_of_week: %d", day)
		}
		seenDays[day] = true
	}
	
	// Validate time format and range
	startTime, err := time.Parse("15:04", rule.StartTime)
	if err != nil {
		return fmt.Errorf("invalid start_time format: %w", err)
	}
	endTime, err := time.Parse("15:04", rule.EndTime)
	if err != nil {
		return fmt.Errorf("invalid end_time format: %w", err)
	}
	if !endTime.After(startTime) {
		return errors.New("end_time must be after start_time")
	}
	
	return nil
}
