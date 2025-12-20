package sqlite

import "database/sql"

type MembershipStore struct {
	db *sql.DB
}
