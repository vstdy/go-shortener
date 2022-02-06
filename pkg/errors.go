package pkg

import "errors"

var (
	ErrNoDBConnection  = errors.New("no DB connection")
	ErrUniqueViolation = errors.New("unique constraint violation error")
)
