package models

import "time"

type Session struct {
	Token     string    `db:"token"`
	UserID    int       `db:"user_id"`
	Login     string    `db:"login"`
	CreatedAt time.Time `db:"created_at"`
}
