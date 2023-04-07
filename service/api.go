package service

import (
	"context"
	"database/sql"

	"github.com/go-kit/kit/log"
)

type service struct {
	db     *sql.DB
	value  string
	logger log.Logger
}

// Service interface describes a service that adds numbers
type Service interface {
	Add(ctx context.Context, numA, numB float32) (float32, error)
}

// NewService returns a Service with all of the expected dependencies
func NewService(db *sql.DB, value string, logger log.Logger) Service {
	return &service{
		db:     db,
		value:  value,
		logger: logger,
	}
}

// Add func implements Service interface
func (s service) Add(ctx context.Context, numA, numB float32) (float32, error) {
	return numA + numB, nil
}
