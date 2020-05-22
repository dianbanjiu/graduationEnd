package controller

import (
	"github.com/gin-gonic/gin"
	"graduationEnd/common"
	"graduationEnd/model"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func AddPublication(ctx *gin.Context) {
	var publication model.Publication
	user, _ := ctx.Get("user")
	if user.(model.User).SelectCourse == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "您还未选择任何课程，请选课之后再添加周报",
		})
		ctx.Abort()
		return
	}
	_ = ctx.Bind(&publication)
	if len(publication.Content) < 10 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "内容不足 10 个字",
		})
		ctx.Abort()
		return
	}

	tn := time.Now()
	t := strings.ReplaceAll(tn.Format("2006-01-02"), "-", "")
	publication.ID = user.(model.User).ID + t
	publication.CreateAt = tn
	publication.StudentID = user.(model.User).ID
	publication.CourseID = user.(model.User).SelectCourse
	publication.StudentName = user.(model.User).Name

	db := common.GetDB()
	var course model.Course
	db.Find(&course, "id = ?", user.(model.User).SelectCourse)
	publication.TeacherID = course.Mentor
	var teacher model.User
	db.Find(&teacher, "id = ?", publication.TeacherID)
	publication.TeacherName = teacher.Name
	db.Save(&publication)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "发布成功",
	})
}

// 查看所有的周报
func ViewAllPublication(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	identify := user.(model.User).Identify
	var publications []model.Publication
	db := common.GetDB()
	if identify == "student" {
		db.Find(&publications, map[string]interface{}{"student_id": user.(model.User).ID, "course_id": user.(model.User).SelectCourse})
	} else if identify == "teacher" {
		db.Find(&publications, "teacher_id = ?", user.(model.User).ID)
	}

	for i := 0; i < len(publications)/2; i++ {
		publications[i], publications[len(publications)-i-1] = publications[len(publications)-i-1], publications[i]
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  publications,
	})
}

// 教师为周报打分
func EvaluationAndScore(ctx *gin.Context) {
	var publication model.Publication
	_ = ctx.Bind(&publication)

	if publication.ID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "请选择周报后进行操作",
		})
		ctx.Abort()
		return
	}
	m, _ := regexp.Match("^1?[1-9]?[0-9]", []byte(publication.TeacherScore))
	if !m {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "格式有误。请注意分数为百分制且为整数，如 89",
		})
		ctx.Abort()
		return
	}

	db := common.GetDB()
	db.Model(&publication).Updates(model.Publication{TeacherEvaluation: publication.TeacherEvaluation, TeacherScore: publication.TeacherScore})
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "评价/评分成功",
	})
}

// 教师删除评价
func DeleteEvaluation(ctx *gin.Context) {
	var publication model.Publication
	_ = ctx.Bind(&publication)
	if publication.ID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "请选择周报后进行操作",
		})
		ctx.Abort()
		return
	}

	db := common.GetDB()
	db.Model(&publication).Updates(map[string]string{"teacher_evaluation": "", "teacher_score": ""})
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "删除成功",
	})
}

func StudentScoreTimes(ctx *gin.Context) {
	studentId := ctx.Query("student_id")
	if len(studentId) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "401",
			"msg":  "学生ID不可为空",
		})
		ctx.Abort()
		return
	}
	db := common.GetDB()
	var publications []model.Publication
	db.Find(&publications, "student_id = ?", studentId)
	var studentScore = make(map[string]int)
	for _, publication := range publications {
		if len(publication.TeacherScore) == 0 {
			continue
		}
		studentScore[publication.TeacherScore] += 1
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  studentScore,
	})
}

// 学生查询实训申请进度
func ViewApplyProgress(ctx *gin.Context) {
	student, _ := ctx.Get("user")
	db := common.GetDB()
	var suggestCourse model.SuggetsCourse
	db.Not("teacher_status", "2").Not("admin_status", 2).Find(&suggestCourse, "student_id = ?", student.(model.User).ID)
	if suggestCourse.CourseName=="" {
		db.Find(&suggestCourse, "student_id = ?", student.(model.User).ID)
	}
	if suggestCourse.ID=="" {
		ctx.JSON(http.StatusNoContent, nil)
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":          "200",
		"msg": suggestCourse,
	})
}

// 学生自主申请实训
func ApplyCourse(ctx *gin.Context){
	student, _ := ctx.Get("user")
	db:=common.GetDB()
	var suggestCourse model.SuggetsCourse
	db.Not("teacher_status", "2").Not("admin_status", 2).Find(&suggestCourse, "student_id = ?", student.(model.User).ID)
	if suggestCourse.CourseName!="" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":"400",
			"msg":"尚有课程在申请中，暂不可再次提交申请",
		})
		ctx.Abort()
		return
	}
    var selectedCourse model.User
	db.Find(&selectedCourse, "id = ?", student.(model.User).ID)
	if selectedCourse.SelectCourse!="" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":"400",
			"msg":"已经选课，暂不可提交申请",
		})
		ctx.Abort()
		return
	}
	_ = ctx.Bind(&suggestCourse)
	if suggestCourse.CourseName==""||suggestCourse.MentorID==""||
		suggestCourse.Address==""||suggestCourse.Company==""||suggestCourse.Desc==""||
		suggestCourse.MentorName==""{
		ctx.JSON(http.StatusBadRequest,gin.H{
			"code":"400",
			"msg":"数据不完整，请修改后重新提交",
		})
		ctx.Abort()
		return
	}

	suggestCourse.StudentID = student.(model.User).ID
	suggestCourse.StudentName = student.(model.User).Name
	suggestCourse.ID=student.(model.User).ID+time.Now().Format("20060102")
	suggestCourse.TeacherStatus="0"
	suggestCourse.AdminStatus="0"
	suggestCourse.CreateAt=time.Now()
	db.Create(&suggestCourse)
	ctx.JSON(http.StatusOK, gin.H{
		"code":"200",
		"msg":"提交成功",
	})
}

// 教师/管理员查看学生的实训申请
func ViewAllApplyProgress(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	var suggestCourses []model.SuggetsCourse
	db := common.GetDB()
	if user.(model.User).Identify=="teacher" {
		db.Find(&suggestCourses, "mentor_id = ?", user.(model.User).ID)
	}else if user.(model.User).Identify == "admin" {
		db.Not("teacher_status", "0").Find(&suggestCourses)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":"200",
		"msg":suggestCourses,
	})
}

// 教师/管理员处理学生的实训申请
func HandleApplyCourse(ctx *gin.Context){
    user, _ := ctx.Get("user")
    db := common.GetDB()
    var suggestCourse model.SuggetsCourse
	_ = ctx.Bind(&suggestCourse)
	if suggestCourse.ID=="" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":"400",
			"msg":"请检查提交项是否完整",
		})
		ctx.Abort()
		return
	}
	db.Model(&model.SuggetsCourse{}).Updates(suggestCourse)
	if user.(model.User).Identify == "admin" && suggestCourse.AdminStatus == "1" {
		courseId := ctx.Query("courseId")
		var course = model.Course{
			ID:       courseId,
			Name:     suggestCourse.CourseName,
			Address:  suggestCourse.Address,
			Desc:     suggestCourse.Desc,
			Company:  suggestCourse.Company,
			Mentor:   suggestCourse.MentorID,
			Count:    1,
			CreateAt: time.Now(),
		}
		var tempCourse model.Course
		db.Find(&tempCourse, "id= ?", courseId)
		if tempCourse.Name=="" {
			db.Save(&course)
		}else {
			db.Model(&course).Update("count",tempCourse.Count+1)
		}

		var student model.User
		db.Find(&student, "id = ?", suggestCourse.StudentID)
		if student.ID!="" {
			db.Model(&student).Update("select_course", courseId)
		}
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":"200",
		"msg":"处理成功",
	})
}