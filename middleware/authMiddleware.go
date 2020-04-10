package middleware

import (
	"github.com/gin-gonic/gin"
	"graduationEnd/common"
	"graduationEnd/model"
	"net/http"
	"strings"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":  "401",
				"msg":   "认证失败",
				"token": "",
			})
			ctx.Abort()
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := common.CheckToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":  "401",
				"msg":   "认证失败",
				"token": "",
			})
			ctx.Abort()
			return
		}

		userID := claims.UserID
		db := common.GetDB()
		var user model.User
		db.First(&user, userID)
		if user.ID == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":  "401",
				"msg":   "认证失败",
				"token": "",
			})
			ctx.Abort()
			return
		}
		ctx.Set("user", user)
		ctx.Next()
	}
}
