package common

import (
	"graduationEnd/model"
)

type UserDto struct {
	Identify   string `form:"identify" xml:"identify" json:"identify"`       //身份
	ID         string `form:"id" xml:"id" json:"id"`                         // 学号/工号
	Name       string `form:"name" xml:"name" json:"name"`                   // 姓名
	Phone      string `form:"phone" xml:"phone" json:"phone"`                //手机号
	School     string `form:"school" xml:"school" json:"school"`             // 学校名/公司名
	Profession string `form:"profession" xml:"profession" json:"profession"` // 专业/职务
	Gender     string `form:"gender" xml:"gender" json:"gender"`             // 性别
}

type BoardDto struct {
	ID       string
	Content  string
	CreateAt string
}

type CourseDto struct {
	ID            string `xml:"id" json:"id" form:"id"`                // 岗位的 ID
	Name          string `xml:"name" json:"name" form:"name"`          // 岗位名
	Address       string `xml:"address" json:"address" form:"address"` //工作地点
	Desc          string `xml:"desc" json:"desc" form:"desc"`          // 岗位的描述
	Company       string `xml:"company" json:"company" form:"company"` // 发布岗位的公司
	Mentor        string `xml:"mentor" json:"mentor" form:"mentor"`    // 对应的指导教师
	Count         int    `xml:"count" json:"count" form:"count"`       // 岗位招收人数
	AlreadySelect int    `xml:"already_select" json:"already_select" form:"already_select"`
	CreateAt      string
}

// 统一用户返回格式
func UserToDto(user model.User) UserDto {
	return UserDto{
		Identify:   user.Identify,
		ID:         user.ID,
		Name:       user.Name,
		Phone:      user.Phone,
		School:     user.School,
		Profession: user.Profession,
		Gender:     user.Gender,
	}
}
func BoardToDto(board model.Board) BoardDto {
	return BoardDto{
		ID:       board.ID,
		Content:  board.Content,
		CreateAt: board.CreateAt.Format("2006-01-02"),
	}
}

func CourseToDto(course model.Course) CourseDto {
	var count int
	db := GetDB()
	db.Model(&model.User{}).Where("select_course= ?", course.ID).Count(&count)
	var teacher model.User
	db.Find(&teacher, "id = ?",course.Mentor)
	return CourseDto{
		ID:            course.ID,
		Name:          course.Name,
		Address:       course.Address,
		Desc:          course.Desc,
		Company:       course.Company,
		Mentor:        teacher.Name,
		Count:         course.Count,
		AlreadySelect: count,
		CreateAt:      course.CreateAt.Format("2006-01-02"),
	}
}
