package pkg

import "errors"

var (
	ErrUnsupportedStorageType = errors.New("unsupported storage type")
	ErrInvalidInput           = errors.New("invalid input")
	ErrAlreadyExists          = errors.New("object exists in the DB")
	ErrNoDBConnection         = errors.New("no DB connection")
)
