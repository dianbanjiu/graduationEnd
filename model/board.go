package model

import (
	"time"
)

type Board struct {
	ID        string    `form:"id" json:"id" xml:"id" gorm:"primary_key"`
	Content   string    `form:"content" json:"content" xml:"content"`
	CreateAt  time.Time `form:"create_at" json:"create_at" xml:"create_at"`
	DeletedAt *time.Time
}
