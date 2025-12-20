package domain

type Membership struct {
	user_id   string
	family_id string
	role      string // (owner | member)
}
