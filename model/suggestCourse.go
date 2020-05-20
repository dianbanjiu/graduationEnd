package model

import "time"

type SuggetsCourse struct {
	ID            string `json:"id" xml:"id" form:"id"`
	CourseName          string `xml:"course_name" json:"course_name" form:"course_name"`          // 岗位名
	Address       string `xml:"address" json:"address" form:"address"` //工作地点
	Desc          string `xml:"desc" json:"desc" form:"desc"`          // 岗位的描述
	Company       string `xml:"company" json:"company" form:"company"` // 发布岗位的公司
	StudentName string `xml:"student_name" json:"student_name" form:"student_name"`
	StudentID string `xml:"student_id" json:"student_id" form:"student_id"`
	MentorID        string `xml:"mentor_id" json:"mentor_id" form:"mentor_id"`
	MentorName    string `json:"mentor_name" xml:"mentor_name" form:"mentor_name"`
	TeacherStatus string `json:"teacher_status" xml:"teacher_status" form:"teacher_status"` // 教师的处理进度，0：未处理，1：同意申请，2：拒绝申请
	AdminStatus   string `json:"admin_status" xml:"admin_status" form:"admin_status"`       // 管理员处理进度，同教师状态码
	CreateAt      time.Time `json:"create_at" xml:"create_at" form:"create_at"`
}
