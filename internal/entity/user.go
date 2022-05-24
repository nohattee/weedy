package entity

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string       `sql:"user_id"`
	FirstName string       `sql:"first_name"`
	LastName  string       `sql:"last_name"`
	Email     string       `sql:"email"`
	CreatedAt time.Time    `sql:"created_at"`
	UpdatedAt time.Time    `sql:"updated_at"`
	DeletedAt sql.NullTime `sql:"deleted_at"`
}
