package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"scheduler-service/internal/models"
)

// Querier abstracts pgx pool/tx for easier testing and transactions.
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type AvailabilityRepository interface {
	InsertAvailabilityRule(ctx context.Context, q Querier, r *models.AvailabilityRule) error
	ListAvailabilityRules(ctx context.Context, q Querier, userID string) ([]models.AvailabilityRule, error)
	ListAvailabilityRulesByEvent(ctx context.Context, q Querier, userID, event string) ([]models.AvailabilityRule, error)
	UpdateAvailabilityRule(ctx context.Context, q Querier, userID, ruleID string, r *models.AvailabilityRule) (string, error)
	GetAvailabilityRule(ctx context.Context, q Querier, userID, ruleID string) (*models.AvailabilityRule, error)
	DeleteAvailabilityRule(ctx context.Context, q Querier, userID, ruleID string) error
}
