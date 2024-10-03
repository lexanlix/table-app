package db

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	ErrRequiredHostPort         = errors.New("invalid db configuration: host and port are required")
	ErrRequiredDatabase         = errors.New("invalid db configuration: database is required")
	ErrRequiredUsernamePassword = errors.New("invalid db configuration: username and password are required")
)

type StorageConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Schema   string `json:"schema"`
}

func (cfg StorageConfig) getDsn() (string, error) {
	if len(cfg.Host) == 0 || len(cfg.Port) == 0 {
		return "", ErrRequiredHostPort
	}

	if len(cfg.Database) == 0 {
		return "", ErrRequiredDatabase
	}

	if len(cfg.Username) == 0 || len(cfg.Password) == 0 {
		return "", ErrRequiredUsernamePassword
	}

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	if len(cfg.Schema) != 0 {
		dsn += fmt.Sprintf("?search_path=%s", cfg.Schema)
	}

	return dsn, nil
}
