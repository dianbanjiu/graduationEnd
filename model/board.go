package model

import (
	"time"
)

type Board struct {
	ID             string    `form:"id" json:"id" xml:"id" gorm:"primary_key"`
	Content        string    `form:"content" json:"content" xml:"content"`
	CreateAt       time.Time `form:"create_at" json:"create_at" xml:"create_at"`
	CreateBy       string    `json:"create_by" xml:"create_by" form:"create_by"`
	CreateIdentify string    `json:"create_identify" xml:"create_identify" form:"create_identify"`
	DeletedAt      *time.Time
}
