package models

import "time"

type Content struct {
	Id        int       `json:"id" db:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	RoomId    string    `json:"room_id" db:"room_id"`
	ToUserId  int       `json:"to_uid_id" db:"to_user_id"`
	Content   string    `json:"content" db:"content"`
	ImageUrl  string    `json:"image_url" db:"image_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" db:"deleted_at"`
}
