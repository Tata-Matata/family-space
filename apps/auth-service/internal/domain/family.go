package domain

import (
	"github.com/google/uuid"

	"time"
)

type Family struct {
	id         uuid.UUID
	name       string
	created_at time.Time
}
