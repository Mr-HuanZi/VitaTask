package controller

import (
	"VitaTaskGo/app/logic"
	"VitaTaskGo/app/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FilesController struct {
}

func NewFilesController() *FilesController {
	return &FilesController{}
}

func (*FilesController) UploadFile(ctx *gin.Context) {
	// 获取字段名
	keyName := ctx.DefaultQuery("key", "file")
	// 指定的文件类型
	fileType := ctx.DefaultQuery("type", "all")
	file, err := ctx.FormFile(keyName)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(logic.NewFilesLogic(ctx).UploadFile(file, fileType)))
}
