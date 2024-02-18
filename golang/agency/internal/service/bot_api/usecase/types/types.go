package types

import (
	uuid "github.com/satori/go.uuid"
)

type Advertisement struct {
	ID uuid.UUID

	// metadata
	Name     string
	Topics   []string
	Coverage int
	Budget   int
	Message  string
}
