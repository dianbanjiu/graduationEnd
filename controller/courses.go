package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"graduationEnd/common"
	"graduationEnd/model"
	"net/http"
	"os"
	"path"
	"time"
)

// 获取全部课程信息
func GetCourses(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	var courses = make([]model.Course, 0)
	db := common.GetDB()
	// 判断用户身份，如果是学生和管理员，则返回所有课程，如果是教师则仅返回本人所负责的课程
	if user.(model.User).Identify == "student" || user.(model.User).Identify == "admin" {
		db.Find(&courses)
	} else if user.(model.User).Identify == "teacher" {
		db.Where("mentor = ?", user.(model.User).ID).Find(&courses)
	}

	// 将返回结果统一格式
	var coursesDto = make([]common.CourseDto, len(courses))
	for i, v := range courses {
		coursesDto[i] = common.CourseToDto(v)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  coursesDto,
	})
}

// 删除多个指定课程，传递参数时需要使用json数组的格式，具体格式如下
// [
//		{"id":"000001"},
//		{"id":"000002"}
// ]
func DeleteCourses(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	if user.(model.User).Identify != "admin" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "权限不足",
		})
		ctx.Abort()
		return
	}

	// 获取所要删除的课程
	var course model.Course
	_ = ctx.Bind(&course)

	// 逐个删除所给的课程
	db := common.GetDB()
		if course.ID != "" {
			db.Delete(&course)
		}

	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "删除成功",
	})
}

//添加单个课程
func AddCourse(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	if user.(model.User).Identify != "admin" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "权限不足",
		})
		ctx.Abort()
		return
	}
	var course model.Course
	_ = ctx.Bind(&course)

	// 课程的重要信息不可以为空
	if course.ID == "" || course.Name == "" || course.Company == "" || course.Mentor == "" || course.Address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "添加失败",
		})
		ctx.Abort()
		return
	}

	// 添加课程到数据库
	course.CreateAt = time.Now()
	db := common.GetDB()
	db.Create(&course)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "添加成功",
	})
}

// 从文件批量添加课程
func AddCourses(ctx *gin.Context) {
	//获取文件
	form, _ := ctx.MultipartForm()
	files := form.File["file"]
	db := common.GetDB()

	// 逐个解析文件，并将文件保存到服务器的 /tmp 目录下
	for _, file := range files {
		dst := path.Join("./tmp/c" + file.Filename)
		_ = ctx.SaveUploadedFile(file, dst)
		f, _ := xlsx.OpenFile(dst)
		for _, sheet := range f.Sheets {
			for i := 2; i < len(sheet.Rows); i++ {
				var course model.Course
				t := sheet.Rows[i].Cells
				if len(t) != 7 {
					continue
				}
				db.Where("id = ?", t[1].String()).First(&course)
				if course.ID == "" {
					course.ID = t[0].String()
					if course.ID == "" {
						continue
					}
					course.Name = t[1].String()
					course.Company = t[2].String()
					course.Address = t[3].String()
					course.Count, _ = t[4].Int()
					course.Mentor = t[5].String()
					course.Desc = t[6].String()
					course.CreateAt = time.Now()
					db.Create(&course)
				}
			}
		}
		_ = os.Remove(dst)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "上传成功",
	})
}

// 学生查询已选课程
func AlreadySelectedCourse(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	if user.(model.User).Identify != "student" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "权限不足",
		})
		ctx.Abort()
		return
	}
	db := common.GetDB()
	var student model.User
	db.Find(&student, "id = ?", user.(model.User).ID)
	if student.SelectCourse == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "尚未选择任何课程",
		})
		ctx.Abort()
		return
	}
	var course model.Course
	db.Find(&course, "id = ?", student.SelectCourse)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  common.CourseToDto(course),
	})
}

// 学生退课
func DropCourse(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	if user.(model.User).Identify != "student" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "权限不足",
		})
		ctx.Abort()
		return
	}
	if user.(model.User).SelectCourse == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "您还尚未选择课程",
		})
		ctx.Abort()
		return
	}
	var student model.User
	student = user.(model.User)
	db := common.GetDB()
	db.Model(&student).Update("select_course", "")
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "退课成功",
	})
}

// 学生选课
func SelectCourse(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	if user.(model.User).Identify != "student" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "权限不足",
		})
		ctx.Abort()
		return
	}
	// 每人限选一门，如果选课之前还有其他已选课程，不予通过
	if user.(model.User).SelectCourse != "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": "400",
			"msg":  "每人限选一门，如需选择其他课程，请先退掉当前课程",
		})
		ctx.Abort()
		return
	}
	var student model.User
	student = user.(model.User)

	// 判断需要选择的课程是否存在于数据库
	var course model.Course
	_ = ctx.Bind(&course)
	if course.ID == "" {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code": "400",
			"msg":  "未找到所选课程",
		})
		ctx.Abort()
		return
	}

	// 判断需选课程人数是否已达上限
	db := common.GetDB()
	db.Find(&course, "id = ?", course.ID)
	var selectedCount int
	var selectedMem []model.User
	db.Find(&selectedMem, "select_course = ?", course.ID)
	selectedCount = len(selectedMem)
	if course.Count-selectedCount <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "该实训人数已满",
		})
		ctx.Abort()
		return
	} else {
		db.Model(&student).Update("select_course", course.ID)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "选课成功",
	})
}
