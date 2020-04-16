package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
	"graduationEnd/common"
	"graduationEnd/model"
	"net/http"
	"os"
	"path"
	"regexp"
)

// 获取用户个人信息
func GetUserInfo(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  common.UserToDto(user.(model.User)),
	})
}

// 修改用户个人信息
func ChangeInfo(ctx *gin.Context) {
	var user model.User
	_ = ctx.Bind(&user)
	var tempUser model.User
	db := common.GetDB()
	db.Find(&tempUser, user.ID)
	if tempUser.ID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code": "401",
			"msg":  "修改失败",
		})
		ctx.Abort()
		return
	}

	m, _ := regexp.Match("^1[0-9]{10}$", []byte(user.Phone))
	if m {
		tempUser.Phone = user.Phone
	}
	if user.Password != "" && len(user.Password) >= 6 {
		tempUser.Password = common.AesEncrypt(user.Password)
	}

	db.Model(&tempUser).Updates(map[string]string{"phone": tempUser.Phone, "password": tempUser.Password})

	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "修改成功",
	})
}

// 获取所有教师或者学生的信息
func GetUsersInfo(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	if user.(model.User).Identify != "admin" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "权限不足",
		})
		ctx.Abort()
		return
	}
	identify := ctx.Query("identify")
	if identify != "student" && identify != "teacher" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "用户身份不存在",
		})
		ctx.Abort()
		return
	}

	var users []model.User
	db := common.GetDB()
	db.Find(&users, "identify = ?", identify)
	var userDto = make([]common.UserDto, len(users))
	for k, v := range users {
		userDto[k] = common.UserToDto(v)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  userDto,
	})
}

// 删除教师/学生们的信息,传递参数时需要使用json数组的格式，具体格式如下
// [
//		{"id":"000001"},
//		{"id":"000002"}
// ]
func DeleteUsers(ctx *gin.Context) {
	var user model.User
	_ = ctx.Bind(&user)
	db := common.GetDB()
		if user.ID != "" {
			db.Delete(&user)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "删除成功",
	})
}

// 添加用户
func AddUser(ctx *gin.Context) {
	var user model.User
	_ = ctx.Bind(&user)
	if user.ID == "" || user.Identify == "" || user.Name == "" || len(user.Phone) != 11 || user.Gender == "" || user.Profession == "" || user.School == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "身份信息不完整",
		})
		ctx.Abort()
		return
	}

	user.Password = common.AesEncrypt(user.ID)
	db := common.GetDB()
	db.Create(&user)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "添加成功",
	})
}

// 添加用户们
func AddUsers(ctx *gin.Context) {
	form, _ := ctx.MultipartForm()
	files := form.File["file"]

	db := common.GetDB()
	for _, file := range files {
		dst := path.Join("./tmp/u" + file.Filename)
		_ = ctx.SaveUploadedFile(file, dst)
		f, _ := xlsx.OpenFile(dst)
		for _, sheet := range f.Sheets {
			for i := 4; i < len(sheet.Rows); i++ {
				var user model.User
				t := sheet.Rows[i].Cells
				if len(t) != 7 {
					continue
				}
				db.Where("id = ?", t[1].String()).First(&user)
				if user.ID == "" {
					user.Identify = t[0].String()
					user.ID = t[1].String()
					user.Name = t[2].String()
					user.Gender = t[3].String()
					user.Phone = t[4].String()
					user.School = t[5].String()
					user.Profession = t[6].String()
					user.Password = common.AesEncrypt(user.ID)
					db.Create(&user)
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
