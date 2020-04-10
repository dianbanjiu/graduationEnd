package model

import "time"

type User struct {
	Identify     string `form:"identify" xml:"identify" json:"identify"`                //身份
	ID           string `form:"id" xml:"id" json:"id"`                                  // 学号/工号
	Name         string `form:"name" xml:"name" json:"name"`                            // 姓名
	Phone        string `form:"phone" xml:"phone" json:"phone"`                         //手机号
	Password     string `form:"password" xml:"password" json:"password"`                //密码
	School       string `form:"school" xml:"school" json:"school"`                      // 学校名/公司名
	Profession   string `form:"profession" xml:"profession" json:"profession"`          // 专业/职务
	Gender       string `form:"gender" xml:"gender" json:"gender"`                      // 性别
	SelectCourse string `xml:"select_course" json:"select_course" form:"select_course"` // 学生的选课情况
	DeleteAt     *time.Time
}
