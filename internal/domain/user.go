package domain

import "time"

type User struct {
	ID             int64     `db:"id" json:"id"`
	Email          string    `db:"email" json:"email"`
	Name           string    `db:"name" json:"name"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	HashedPassword string    `db:"hashed_password" json:"-"`
	ProfilePicture string    `db:"profile_picture" json:"profilePicture"`
}

func NewUser(email, name string, hashedPass, profilePic string) *User {
	return &User{
		Email:          email,
		Name:           name,
		CreatedAt:      time.Now(),
		HashedPassword: hashedPass,
		ProfilePicture: profilePic,
	}
}
