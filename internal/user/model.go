package user

import (
	"context"
	"time"
)

// User is a user in the system
type User struct {
	tableName struct{} `pg:"users,alias:users"`

	ID       int    `pg:",pk" json:"-"`
	Email    string `pg:",unique,notnull" json:"email"`
	Mobile   string `pg:",unique,notnull" json:"mobile"`
	Password string `pg:",notnull" json:"-"`

	FirstName string `pg:",notnull" json:"first_name"`
	LastName  string `pg:",notnull" json:"last_name"`
	ImageURL  string `pg:",notnull" json:"image_url"`
	Address   string `pg:",notnull" json:"address"`

	Active bool `pg:",notnull" json:"active"`

	CreatedAt time.Time  `pg:",notnull" json:"created_at"`
	UpdatedAt time.Time  `pg:",notnull" json:"updated_at"`
	DeletedAt *time.Time `pg:",soft_delete" json:"-"`
}

// BeforeInsert Before insert trigger
func (o *User) BeforeInsert(c context.Context) (context.Context, error) {
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()

	return c, nil
}

// BeforeUpdate Before Update trigger
func (o *User) BeforeUpdate(c context.Context) (context.Context, error) {
	o.UpdatedAt = time.Now()

	return c, nil
}
