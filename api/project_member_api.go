package api

import (
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/response"
	"VitaTaskGo/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProjectMemberApi struct {
}

func NewProjectMemberApi() *ProjectMemberApi {
	return &ProjectMemberApi{}
}

func (receiver ProjectMemberApi) List(ctx *gin.Context) {
	var post dto.ProjectMemberListQuery
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	list := service.NewProjectMemberService(db.Db, ctx).GetMembers(post)
	ctx.JSON(http.StatusOK, response.SuccessData(list))
}

func (receiver ProjectMemberApi) Remove(ctx *gin.Context) {
	var post dto.ProjectMemberBind
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(nil, service.NewProjectMemberService(db.Db, ctx).RemoveRole(post.ProjectId, post.UserId, post.Role)))
}

// Bind 新增普通成员
func (receiver ProjectMemberApi) Bind(ctx *gin.Context) {
	var post dto.ProjectMemberBind
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(nil, service.NewProjectMemberService(db.Db, ctx).Bind(post.ProjectId, post.UserId, post.Role)))
}
