package database

import "time"

type Status string

// List of Status
const (
	CREATED       Status = "created"
	APPROVED      Status = "approved"
	DECLINED      Status = "declined"
	ON_MODERATION Status = "on moderation"
)

type Flat struct {
	Id      int    `json:"id" binding:"required"`
	HouseId int    `json:"house_id" binding:"required"`
	Price   int    `json:"price" binding:"required"`
	Rooms   int    `json:"rooms" binding:"required"`
	Status  Status `json:"status" binding:"required"`
}

type House struct {
	Id        int        `json:"id" binding:"required"`
	Address   string     `json:"address" binding:"required"`
	Year      int        `json:"year" binding:"required"`
	Developer string     `json:"developer,omitempty" binding:"required"`
	CreatedAt *time.Time `json:"created_at,omitempty" binding:"required"`
	UpdateAt  *time.Time `json:"update_at,omitempty" binding:"required"`
}

type User struct {
	UserId   string
	Email    string
	Password string
	UserType string
}
