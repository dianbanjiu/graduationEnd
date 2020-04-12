package controller

import (
	"github.com/gin-gonic/gin"
	"graduationEnd/common"
	"graduationEnd/model"
	"net/http"
)

// 登录认证
func LoginAuth(ctx *gin.Context) {
	// 获取用户的输入字段
	var user model.User
	_ = ctx.Bind(&user)
	db := common.GetDB()
	var tempUser model.User
	// 查询数据库中是否有对应的用户字段
	db.First(&tempUser, user.ID)
	// 判断用户名及密码是否正确
	if tempUser.Identify == user.Identify && tempUser.ID == user.ID && common.AesDecrypt(tempUser.Password) == user.Password {
		ctx.JSON(http.StatusOK, gin.H{
			"code":  "200",
			"msg":   "认证成功",
			"token": common.ReleaseToken(user.ID),
		})
	} else {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":  "401",
			"msg":   "认证失败",
			"token": "",
		})
		ctx.Abort()
		return
	}
}
