package pkg

import "errors"

var (
	ErrUnsupportedStorageType = errors.New("unsupported storage type")
	ErrInvalidInput           = errors.New("invalid input")
	ErrNoValue                = errors.New("value is missing")
	ErrNoDBConnection         = errors.New("no DB connection")
	ErrIntegrityViolation     = errors.New("integrity constraint violation error")
)
