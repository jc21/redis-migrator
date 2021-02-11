package model

import "fmt"

// RedisServerConfig ...
type RedisServerConfig struct {
	Hostname string
	Port     int
	DBIndex  int
	Username string
	Password string
}

// Check ensures that the details given are enough
func (c *RedisServerConfig) Check() error {
	// Ensure password is given if username is given
	if c.Username != "" && c.Password == "" {
		return fmt.Errorf("Password is required when username is supplied")
	}
	return nil
}
