package pkg

import "errors"

var (
	ErrNoDBConnection     = errors.New("no DB connection")
	ErrIntegrityViolation = errors.New("integrity constraint violation error")
)
