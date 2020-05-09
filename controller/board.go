package controller

import (
	"github.com/gin-gonic/gin"
	"graduationEnd/common"
	"graduationEnd/model"
	"net/http"
	"strconv"
	"time"
)

// 返回公告内容
func GetBoards(ctx *gin.Context) {
	db := common.GetDB()
	var boards = make([]model.Board, 0)

	userInfo, _ := ctx.Get("user")
	switch userInfo.(model.User).Identify {
	case "admin":
		db.Find(&boards, "create_identify = ?", "admin")
		break
	case "student":
		var selectedCourse model.Course
		db.Find(&selectedCourse, "id = ?", userInfo.(model.User).SelectCourse)
		db.Where("create_by = ?", selectedCourse.Mentor).
			Or("create_identify = ?", "admin").Find(&boards)
		break
	case "teacher":
		db.Where("create_by = ?", userInfo.(model.User).ID).
			Or("create_identify = ?", "admin").Find(&boards)
		break
	default:
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "Identify is invalid",
		})
		ctx.Abort()
		return
	}

	for i := 0; i < len(boards)/2; i++ {
		boards[i], boards[len(boards)-1-i] = boards[len(boards)-i-1], boards[i]
	}
	var boardDto = make([]common.BoardDto, len(boards))
	for i, v := range boards {
		boardDto[i] = common.BoardToDto(v)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  boardDto,
	})

}

// 添加新的公告板内容
func AddBoard(ctx *gin.Context) {
	var board model.Board
	ctx.Bind(&board)
	if len(board.Content) <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "公告内容不可为空",
		})
		ctx.Abort()
		return
	}
	db := common.GetDB()
	var count int
	db.Table("boards").Count(&count)
	board.ID = strconv.Itoa(count + 1)
	for i := len(board.ID); i <= 8-len(board.ID); {
		board.ID = "0" + board.ID
	}
	userInfo, _ := ctx.Get("user")
	board.CreateBy = userInfo.(model.User).ID
	board.CreateIdentify = userInfo.(model.User).Identify
	board.CreateAt = time.Now()
	db.Create(&board)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "添加成功",
	})
}

// 删除公告
func DeleteBoard(ctx *gin.Context) {
	var board model.Board
	ctx.Bind(&board)
	if board.ID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": "400",
			"msg":  "删除失败",
		})
		ctx.Abort()
		return
	}
	db := common.GetDB()
	db.Delete(&board)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"msg":  "删除成功",
	})
}
