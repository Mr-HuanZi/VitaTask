package handle

import (
	"VitaTaskGo/internal/api/service"
	"VitaTaskGo/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FilesApi struct {
}

func NewFilesApi() *FilesApi {
	return &FilesApi{}
}

func (*FilesApi) UploadFile(ctx *gin.Context) {
	// 获取字段名
	keyName := ctx.DefaultQuery("key", "file")
	// 指定的文件类型
	fileType := ctx.DefaultQuery("type", "all")
	file, err := ctx.FormFile(keyName)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(service.NewFilesService(ctx).UploadFile(file, fileType)))
}
