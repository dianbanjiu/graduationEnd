package model

import "time"

type Course struct {
	ID       string    `xml:"id" json:"id" form:"id"`                      // 岗位的 ID
	Name     string    `xml:"name" json:"name" form:"name"`                // 岗位名
	Address  string    `xml:"address" json:"address" form:"address"`       //工作地点
	Desc     string    `xml:"desc" json:"desc" form:"desc" gorm:"type:text"`                // 岗位的描述
	Company  string    `xml:"company" json:"company" form:"company"`       // 发布岗位的公司
	Mentor   string    `xml:"mentor" json:"mentor" form:"mentor"`          // 对应的指导教师
	Count    int       `xml:"count" json:"count" form:"count"`             // 岗位招收人数
	CreateAt time.Time `form:"create_at" json:"create_at" xml:"create_at"` // 岗位发布时间
	DeleteAt *time.Time
}
