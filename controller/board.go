package controller

import (
	"github.com/gin-gonic/gin"
	"graduationEnd/common"
	"graduationEnd/model"
	"net/http"
	"strconv"
	"time"
)

// 返回最近的5条公告内容
func GetBoards(ctx *gin.Context) {
	db := common.GetDB()
	var boards = make([]model.Board, 0)
	db.Order("create_at").Limit(5).Find(&boards)
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
			"msg":  "添加失败",
		})
		ctx.Abort()
		return
	}
	db := common.GetDB()

	var tempBoard model.Board
	db.Last(&tempBoard)
	id, _ := strconv.Atoi(tempBoard.ID)
	board.ID = strconv.Itoa(id + 1)
	for i := len(board.ID); i <= 8-len(board.ID); {
		board.ID = "0" + board.ID
	}
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
