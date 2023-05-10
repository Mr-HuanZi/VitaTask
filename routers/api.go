package routers

import (
	"VitaTaskGo/app/controller"
	"VitaTaskGo/app/middleware"
	"VitaTaskGo/app/modules/ws"
	"github.com/gin-gonic/gin"
)

func ApiRouters(r *gin.Engine) {
	// 静态文件服务
	r.Static("/uploads", "./uploads")
	// API路由z
	{
		var indexController controller.IndexController
		r.GET("/hello", indexController.Hello)
	}

	{
		// 登录接口
		loginController := controller.NewLoginController()
		r.POST("/login", loginController.Login)
		r.POST("/register", loginController.Register)
	}

	{
		userController := controller.NewUserController()
		// 获取当前登录用户
		r.GET("/currentUser", middleware.CheckLogin(), userController.CurrUser)
		// 用户接口
		g := r.Group("user", middleware.CheckLogin())
		g.POST("store", userController.StoreSelf)
		g.POST("change-avatar", userController.ChangeAvatar)
		g.POST("change-pass", userController.ChangePassword)
		g.POST("change-mobile", userController.ChangeMobile)
		g.POST("change-email", userController.ChangeEmail)
	}

	{
		// 文件接口
		filesController := controller.NewFilesController()
		g := r.Group("/files", middleware.CheckLogin())
		g.POST("upload", filesController.UploadFile)
	}

	{
		// 成员接口
		memberController := controller.NewMemberController()
		g := r.Group("/member", middleware.CheckLogin())
		g.POST("list/simple", memberController.SimpleList)
		g.POST("lists", memberController.Lists)
		g.POST("create", memberController.Create)
		g.POST("disable", memberController.Disable)
		g.POST("enable", memberController.Enable)
		g.POST("reset-pass", memberController.ResetPassword)
		g.POST("change-super", memberController.ChangeSuper)
	}

	{
		// 项目接口
		projectController := controller.NewProjectController()
		g := r.Group("/project", middleware.CheckLogin())
		g.POST("create", projectController.CreateProject)
		g.POST("edit", projectController.EditProject)
		g.POST("list", projectController.ProjectList)
		g.POST("list/simple", projectController.SimpleList)
		g.POST("trash", projectController.ProjectTrash)
		g.POST("del", projectController.ProjectDelete)
		g.POST("archive", projectController.ProjectArchive)
		g.POST("un-archive", projectController.UnArchive)
		g.POST("star", projectController.Star)
		g.POST("un-star", projectController.UnStart)
		g.POST("transfer", projectController.Transfer)
		g.POST("detail", projectController.Detail)

		{
			projectMemberController := controller.NewProjectMemberController()
			gg := g.Group("/member")
			gg.POST("bind", projectMemberController.Bind)
			gg.POST("remove", projectMemberController.Remove)
			gg.POST("list", projectMemberController.List)

		}
	}

	{
		// 任务接口
		taskController := controller.NewTaskController()
		g := r.Group("/task", middleware.CheckLogin())
		g.POST("list", taskController.Lists)
		g.POST("create", taskController.Create)
		g.POST("detail", taskController.Detail)
		g.POST("roles", taskController.Roles)
		g.POST("status", taskController.Status)
		g.POST("change-status", taskController.ChangeStatus)
		g.POST("update", taskController.Update)
		g.POST("delete", taskController.Delete)
		g.POST("statistics", taskController.Statistics)
		g.POST("daily-situation", taskController.DailySituation)

		{
			// 任务组接口
			taskGroupController := controller.NewTaskGroupController()
			gg := g.Group("group")
			gg.POST("add", taskGroupController.Add)
			gg.POST("update", taskGroupController.Update)
			gg.POST("delete", taskGroupController.Delete)
			gg.POST("list", taskGroupController.List)
			gg.POST("detail", taskGroupController.Detail)
			gg.POST("simple-list", taskGroupController.SimpleList)
		}

		{
			taskLogController := controller.NewTaskLogController()
			gg := g.Group("log")
			gg.POST("list", taskLogController.List)
			gg.POST("operators", taskLogController.Operators)
		}
	}

	{
		dialogController := controller.NewDialogController()
		g := r.Group("dialog", middleware.CheckLogin())
		g.POST("list", dialogController.List)
		g.POST("create", dialogController.Create)
		g.POST("msg-list", dialogController.MsgList)
		g.POST("send-text", dialogController.SendText)
	}
}

func WebSocketRouters(r *gin.Engine) {
	// 聊天WS，需要校验登录
	r.GET("chat/", middleware.CheckLogin(), func(ctx *gin.Context) {
		ws.ClientHandle(ctx)
	})
}
