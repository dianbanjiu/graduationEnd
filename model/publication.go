package model

import "time"

type Publication struct{
	ID string	`json:"id" xml:"id" form:"id"`//日志的编号，由学生学号+日期组成，如 1605111420200401
	StudentID	string	`json:"student_id" xml:"student_id" form:"student_id"`
	StudentName string	`json:"student_name" xml:"student_name" form:"student_name"`
	Content string	`json:"content" xml:"content" form:"content"`
	CreateAt time.Time	`json:"create_at" xml:"create_at" form:"create_at"`
	CourseID string `json:"course_id" xml:"course_id" form:"course_id"`
	TeacherID string `json:"teacher_id" xml:"teacher_id" form:"teacher_id"`
	TeacherName string 	`json:"teacher_name" xml:"teacher_name" form:"teacher_name"`
	TeacherEvaluation string `json:"teacher_evaluation" xml:"teacher_evaluation" form:"teacher_evaluation"`
	TeacherScore string	`json:"teacher_score" xml:"teacher_score" form:"teacher_score"`
}