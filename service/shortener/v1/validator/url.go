package validator

import (
	"net/url"
)

// ValidateURL validates email.
func ValidateURL(urlRaw string) error {
	_, err := url.ParseRequestURI(urlRaw)

	return err
}
