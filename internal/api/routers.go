package api

import (
	"VitaTaskGo/internal/api/handle"
	"VitaTaskGo/internal/api/middleware"
	"VitaTaskGo/internal/pkg/ws"
	"github.com/gin-gonic/gin"
)

func Routers(r *gin.Engine) {
	// 静态文件服务
	r.Static("/uploads", "./uploads")

	{
		// 登录接口
		loginApi := handle.NewLoginApi()
		r.POST("/login", loginApi.Login)
		r.POST("/register", loginApi.Register)
	}

	{
		userApi := handle.NewUserApi()
		// 获取当前登录用户
		r.GET("/currentUser", middleware.CheckLogin(), userApi.CurrUser)
		// 用户接口
		g := r.Group("user", middleware.CheckLogin())
		g.POST("store", userApi.StoreSelf)
		g.POST("change-avatar", userApi.ChangeAvatar)
		g.POST("change-pass", userApi.ChangePassword)
		g.POST("change-mobile", userApi.ChangeMobile)
		g.POST("change-email", userApi.ChangeEmail)
	}

	{
		// 文件接口
		filesApi := handle.NewFilesApi()
		g := r.Group("/files", middleware.CheckLogin())
		g.POST("upload", filesApi.UploadFile)
	}

	{
		// 成员接口
		memberApi := handle.NewMemberApi()
		g := r.Group("/member", middleware.CheckLogin())
		g.POST("list/simple", memberApi.SimpleList)
		g.POST("lists", memberApi.Lists)
		g.POST("create", memberApi.Create)
		g.POST("disable", memberApi.Disable)
		g.POST("enable", memberApi.Enable)
		g.POST("reset-pass", memberApi.ResetPassword)
		g.POST("change-super", memberApi.ChangeSuper)
	}

	{
		// 项目接口
		projectApi := handle.NewProjectApi()
		g := r.Group("/project", middleware.CheckLogin())
		g.POST("create", projectApi.CreateProject)
		g.POST("edit", projectApi.EditProject)
		g.POST("list", projectApi.ProjectList)
		g.POST("list/simple", projectApi.SimpleList)
		g.POST("trash", projectApi.ProjectTrash)
		g.POST("del", projectApi.ProjectDelete)
		g.POST("archive", projectApi.ProjectArchive)
		g.POST("un-archive", projectApi.UnArchive)
		g.POST("star", projectApi.Star)
		g.POST("un-star", projectApi.UnStart)
		g.POST("transfer", projectApi.Transfer)
		g.POST("detail", projectApi.Detail)

		{
			projectMemberApi := handle.NewProjectMemberApi()
			gg := g.Group("/member")
			gg.POST("bind", projectMemberApi.Bind)
			gg.POST("remove", projectMemberApi.Remove)
			gg.POST("list", projectMemberApi.List)

		}
	}

	{
		// 任务接口
		taskApi := handle.NewTaskApi()
		g := r.Group("/task", middleware.CheckLogin())
		g.POST("list", taskApi.Lists)
		g.POST("create", taskApi.Create)
		g.POST("detail", taskApi.Detail)
		g.POST("roles", taskApi.Roles)
		g.POST("status", taskApi.Status)
		g.POST("change-status", taskApi.ChangeStatus)
		g.POST("update", taskApi.Update)
		g.POST("delete", taskApi.Delete)
		g.POST("statistics", taskApi.Statistics)
		g.POST("daily-situation", taskApi.DailySituation)

		{
			// 任务组接口
			taskGroupApi := handle.NewTaskGroupApi()
			gg := g.Group("group")
			gg.POST("add", taskGroupApi.Add)
			gg.POST("update", taskGroupApi.Update)
			gg.POST("delete", taskGroupApi.Delete)
			gg.POST("list", taskGroupApi.List)
			gg.POST("detail", taskGroupApi.Detail)
			gg.POST("simple-list", taskGroupApi.SimpleList)
		}

		{
			taskLogApi := handle.NewTaskLogApi()
			gg := g.Group("log")
			gg.POST("list", taskLogApi.List)
			gg.POST("operators", taskLogApi.Operators)
		}
	}

	{
		dialogApi := handle.NewDialogApi()
		g := r.Group("dialog", middleware.CheckLogin())
		g.POST("create", dialogApi.Create)
		g.POST("msg-list", dialogApi.MsgList)
		g.POST("send-text", dialogApi.SendText)
	}

	{
		workflowApi := handle.NewWorkflowApi()
		g := r.Group("workflow", middleware.CheckLogin())
		g.POST("initiate", workflowApi.Initiate)
		g.POST("examine-approve", workflowApi.ExamineApprove)
		g.POST("all", workflowApi.All)
		g.POST("todo", workflowApi.ToDo)
		g.POST("handled", workflowApi.Handled)
		g.POST("list", workflowApi.List)
		g.GET("status/list", workflowApi.StatusList)

		{
			twoG := g.Group("type")
			twoG.POST("add", workflowApi.TypeAdd)
			twoG.POST("update", workflowApi.TypeUpdate)
			twoG.POST("list", workflowApi.TypeList)
			twoG.POST("delete", workflowApi.TypeDelete)
			twoG.POST("detail", workflowApi.TypeDetail)
			twoG.GET("options", workflowApi.TypeOptions)
		}

		{
			twoG := g.Group("node")
			twoG.POST("add", workflowApi.NodeAdd)
			twoG.POST("update", workflowApi.NodeUpdate)
			twoG.POST("list", workflowApi.NodeList)
			twoG.POST("delete", workflowApi.NodeDelete)
			twoG.POST("actions", workflowApi.Actions)
		}
	}
}

func WebSocketRouters(r *gin.Engine) {
	// 聊天WS，需要校验登录
	r.GET("chat/", middleware.CheckLogin(), func(ctx *gin.Context) {
		ws.ClientHandle(ctx)
	})
}
