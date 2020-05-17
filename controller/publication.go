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
		publications[i], publications[len(publications)-i-1] = publications[len(publications)-i-1],publications[i]
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
