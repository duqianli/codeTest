package router

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "service/docs"
	"service/middleWares"
	"service/service"
)

func Router() *gin.Engine {
	ginServer := gin.Default()
	//公有方法
	//路由配置
	ginServer.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	ginServer.GET("/ping", service.Ping)
	//问题
	ginServer.GET("/problem_list", service.GetProblemList)
	ginServer.GET("/problem_detail", service.GetProblemDetail)

	//用户
	ginServer.GET("/user_detail", service.GetUserDetail)

	//提交记录
	ginServer.GET("/submit_list", service.GetSubmitList)
	//排行榜
	ginServer.GET("/rank_list", service.GetRankList)
	//登录
	ginServer.POST("/login", service.Login)
	//注册
	ginServer.POST("/register", service.Register)
	//发送验证码
	ginServer.POST("/sendCode", service.SendCode)

	//管理员私有方法
	authAdmin := ginServer.Group("/admin", middleWares.AuthAdminCheck())
	//问题创建
	authAdmin.POST("/problem-create", service.PloblemCreate)
	//问题修改
	authAdmin.PUT("/problem-modify", service.ProblemModify)
	//获取分类列表
	authAdmin.GET("/Category_list", service.CategoryList)
	//分类列表的创建
	authAdmin.POST("/Category_create", service.CategoryCreate)
	//分类列表的修改
	authAdmin.PUT("/Category_modify", service.CategoryModify)
	//分类列表的删除
	authAdmin.DELETE("/Category_delete", service.CategoryDelete)

	//用户私有方法
	authUser := ginServer.Group("/user", middleWares.AuthUserCheck())
	authUser.POST("/submit", service.Submit)
	return ginServer
}
