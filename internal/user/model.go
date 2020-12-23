package user

import (
	"context"
	"time"
)

// User is a user in the system
type User struct {
	tableName struct{} `pg:"users,alias:users"`

	ID       int    `pg:",pk" json:"-"`
	Email    string `pg:",unique,notnull,use_zero" json:"email"`
	Mobile   string `pg:",unique,notnull,use_zero" json:"mobile"`
	Password string `pg:",notnull" json:"-"`

	FirstName string `pg:",notnull,use_zero" json:"first_name"`
	LastName  string `pg:",notnull,use_zero" json:"last_name"`
	ImageURL  string `pg:",notnull,use_zero" json:"image_url"`
	Address   string `pg:",notnull,use_zero" json:"address"`

	Active bool `pg:",notnull" json:"active"`

	CreatedAt time.Time  `pg:",notnull,use_zero" json:"created_at"`
	UpdatedAt time.Time  `pg:",notnull,use_zero" json:"updated_at"`
	DeletedAt *time.Time `pg:",soft_delete" json:"-"`
}

// Update is a user update in the system
type Update struct {
	Email     *string
	Mobile    *string
	Password  *string
	FirstName *string
	LastName  *string
	ImageURL  *string
	Address   *string
	Active    *bool
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

func populateUpdateInUser(u *User, update *Update) *User {
	if update.Active != nil {
		u.Active = *update.Active
	}

	if update.Email != nil {
		u.Email = *update.Email
	}

	if update.Mobile != nil {
		u.Mobile = *update.Mobile
	}

	if update.Password != nil {
		u.Password = *update.Password
	}

	if update.FirstName != nil {
		u.FirstName = *update.FirstName
	}

	if update.LastName != nil {
		u.LastName = *update.LastName
	}

	if update.ImageURL != nil {
		u.ImageURL = *update.ImageURL
	}

	if update.Address != nil {
		u.Address = *update.Address
	}

	return u
}
