package models

import "time"

type Document struct {
	ID         string    `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Mime       string    `json:"mime" db:"mime"`
	File       bool      `json:"file" db:"file"`
	Public     bool      `json:"public" db:"public"`
	OwnerLogin string    `json:"owner_login" db:"owner_login"`
	Grant      []string  `json:"grant" db:"grant_list"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	JSONData   []byte    `json:"json_data" db:"json_data"`
	FilePath   string    `json:"file_path" db:"file_path"`
}
