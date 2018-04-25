package config

import (
	"fmt"
)

// Config ...
type Config struct {
	DB      DB
	Service Service
}

// DB ...
type DB struct {
	User     string
	Password string
	Database string
	Host     string
	Port     string
}

// Service ...
type Service struct {
	Name string
	Host string
	Port string
}

// Addr ...
func (db DB) Addr() string {
	return fmt.Sprintf("%s:%s", db.Host, db.Port)
}

// Addr ...
func (s Service) Addr() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}
