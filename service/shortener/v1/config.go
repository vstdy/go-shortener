package shortener

import (
	"fmt"
	"time"
)

// Config keeps Service params.
type Config struct {
	DelReqTimeout     time.Duration `mapstructure:"del_req_timeout"`
	DelBufWipeTimeout time.Duration `mapstructure:"del_buf_wipe_timeout"`
	DelBufCap         int           `mapstructure:"del_buf_cap"`
}

// Validate performs a basic validation.
func (config Config) Validate() error {
	if config.DelReqTimeout < time.Second {
		return fmt.Errorf("%s field: too short period", "del_req_timeout")
	}

	if config.DelBufWipeTimeout < time.Second {
		return fmt.Errorf("%s field: too short period", "del_buf_wipe_timeout")
	}

	if config.DelBufCap < 1 {
		return fmt.Errorf("%s field: too small value", "del_buf_cap")
	}

	return nil
}

// NewDefaultConfig builds a Config with default values.
func NewDefaultConfig() Config {
	return Config{
		DelReqTimeout:     5 * time.Second,
		DelBufWipeTimeout: 10 * time.Second,
		DelBufCap:         10,
	}
}
