package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// 跨域问题解决
func Cros() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:9091")
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if ctx.Request.Method == http.MethodOptions {
			ctx.AbortWithStatus(http.StatusOK)
		} else {
			ctx.Next()
		}
	}

}
