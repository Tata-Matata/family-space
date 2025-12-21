package domain

import "time"

type Membership struct {
	UserID    string
	FamilyID  string
	Role      string // (owner | member)
	CreatedAt time.Time
}
