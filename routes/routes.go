package routes

import (
	"fmt"
	cache "github.com/chenyahui/gin-cache"
	"github.com/chenyahui/gin-cache/persist"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"ruoyi/controllers"
	"ruoyi/logger"
	"ruoyi/middlewares"
	"time"
)

func SetUp(mode, host string, port int) *gin.Engine {

	redisStore := persist.NewRedisStore(redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    fmt.Sprintf("%v:%v", host, port),
	}))
	handlerFunc := cache.CacheByRequestURI(redisStore, 2*time.Second)
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode) //gin设置成发布模式
	}
	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"https:holoforever.fun"},                   //允许跨域来发请求的网站
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, //允许的请求方法
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool { //自定义过滤源的方法
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	r.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/captchaImage", controllers.CaptchaImageHandler) //生成二维码，有
	r.POST("/login", controllers.LoginHandler)

	v1 := r.Group("/")
	{
		auth := v1.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("getInfo", handlerFunc, controllers.GetInfoHandler)
			auth.POST("logout", controllers.LogoutHandler)
			auth.GET("getRouters", handlerFunc, controllers.GetRoutersHandler)
			auth.GET("send", controllers.SendEmail)
			auth.GET("/yanzzheng/:id", controllers.YanzhengEmail)
		}
	}
	v2 := r.Group("system/user")
	{
		auth := v2.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.PUT("profile/updatePwd", controllers.UpdatePwdHandler)
			auth.GET("profile", controllers.ProfileHandler)
			auth.PUT("profile", controllers.PostProfileHandler)
			auth.POST("profile/avatar", controllers.AvatarHandler)
		}
	}

	v3 := r.Group("system/dict")
	{
		auth := v3.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("/data/list", controllers.ListDict)
			auth.POST("/data/export", controllers.ExportDict)
			auth.GET("/data/:dictCode", controllers.GetDictCode)
			auth.GET("/data/type/:dictType", controllers.DictTypeHandler)
			auth.POST("/data", controllers.SaveDictData)
			auth.PUT("/data", controllers.UpDictData)
			auth.DELETE("/data/:dictCodes", controllers.DeleteDictData)

			//字典类型-字典管理
			auth.GET("/type/list", controllers.ListDictType)
			auth.POST("/type/export", controllers.ExportType)
			auth.GET("/type/:dictId", controllers.GetTypeDict)
			auth.POST("/type", controllers.SaveType)
			auth.PUT("/type", controllers.UploadType)
			auth.DELETE("/type/:dictIds", controllers.DeleteDataType)
			auth.DELETE("/refreshCache", controllers.RefreshCache)
			auth.GET("/type/optionselect", controllers.GetOptionSelect)
		}
	}
	v4 := r.Group("system/zero")
	{
		auth := v4.Group("")
		{
			auth.GET("/user/list", controllers.ListUser)
			auth.POST("/user/export", controllers.ExportExport)
			auth.POST("/user/importData", controllers.ImportUserData)
			auth.POST("/user/importTemplate", controllers.ImportTemplate)
			auth.GET("/user/:userId", controllers.GetUserInfo)
			auth.GET("/user/", controllers.GetUserInfo)
			auth.POST("/user", controllers.SaveUser)
			auth.PUT("/user", controllers.UploadUser)
			auth.DELETE("/user/:userIds", controllers.DeleteUserById)
			auth.PUT("/user/resetPwd", controllers.ResetPwd)
			auth.PUT("/user/changeStatus", controllers.ChangeUserStatus)
			auth.GET("/user/authRole/:userId", controllers.GetAuthUserRole)
			auth.PUT("/user/authRole", controllers.PutAuthUser)
			auth.GET("/user/deptTree", controllers.GetUserDeptTree)

		}
	}
	v5 := r.Group("system/menu")
	{
		auth := v5.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("/list", controllers.ListMenu)
			auth.GET("/:menuId", controllers.GetMenuInfo)
			auth.GET("/treeselect", controllers.GetTreeSelect)
			auth.GET("/roleMenuTreeselect/:roleId", controllers.TreeSelectByRole)
			auth.POST("", controllers.SaveMenu)
			auth.PUT("", controllers.UploadMenu)
			auth.DELETE("/:menuId", controllers.DeleteMenu)
		}
	}

	v6 := r.Group("system/post")
	{
		auth := v6.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("/list", controllers.ListPost)
			auth.POST("/export", controllers.ExportPost)
			auth.GET("/:postId", controllers.GetPostInfo)
			auth.POST("", controllers.SavePost)
			auth.PUT("", controllers.UploadPost)
			auth.DELETE("/:postIds", controllers.DeletePost)
			auth.GET("/optionselect", controllers.GetPostOptionSelect)
		}
	}
	v7 := r.Group("system/notice")
	{
		auth := v7.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("/list", controllers.ListNotice)
			auth.GET("/:noticeId", controllers.GetNotice)
			auth.POST("", controllers.SaveNotice)
			auth.PUT("", controllers.UploadNotice)
			auth.DELETE("", controllers.DeleteNotice)
		}
	}
	v8 := r.Group("system/role")
	{
		auth := v8.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("/list", controllers.ListRole)
			auth.POST("/export", controllers.ExportRole)
			auth.GET("/:roleId", controllers.GetRoleInfo)
			auth.POST("", controllers.SaveRole)
			auth.PUT("", controllers.UploadRole)
			auth.DELETE("/:roleIds", controllers.DeleteRole) //传入的数据为1,2,3这种类型
			auth.GET("/optionselect", controllers.GetRoleOptionSelect)
			auth.GET("/authUser/allocatedList", controllers.GetAllocatedList)
			auth.GET("/authUser/unallocatedList", controllers.GetUnAllocatedList)
			auth.PUT("/authUser/cancel", controllers.CancelRole)
			auth.PUT("/authUser/cancelAll", controllers.CancelAllRole)
			auth.PUT("/authUser/selectAll", controllers.SelectRoleAll)
			auth.GET("/deptTree/:roleId", controllers.GetDeptTreeRole)

		}
	}
	v9 := r.Group("system/config")
	{
		auth := v9.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("/list", controllers.ListConfig)
			auth.POST("/export", controllers.ExportConfig)
			auth.GET("/:configId", controllers.GetConfigInfo)
			auth.GET("/configKey/:configKey", controllers.GetConfigKey)
			auth.POST("", controllers.SaveConfig)
			auth.PUT("", controllers.UploadConfig)
			auth.DELETE("/:configIds", controllers.DeletectConfig)
			auth.DELETE("/donws/:refreshCache", controllers.DeleteCacheConfig)

		}
	}

	v10 := r.Group("system/dept")
	{
		auth := v10.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("/list", controllers.ListDept)
			auth.GET("/list/exclude/:deptId", controllers.ExcludeDept)
			auth.GET("/:deptId", controllers.GetDept)
			auth.POST("", controllers.SaveDept)
			auth.PUT("", controllers.UpDataDept)
			auth.DELETE("/:deptId", controllers.DeleteDept)
		}
	}
	v11 := r.Group("monitor/cache")
	{
		auth := v11.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("", controllers.CacheHandler)
			auth.GET("getNames", controllers.CacheHandler)
			auth.GET("getKeys/:cacheName", controllers.GetCacheKeysHandler)
			auth.GET("getValue/:cacheName/:cacheKey", controllers.GetCacheValueHandler)
			auth.DELETE("clearCacheName/:cacheName", controllers.ClearCacheNameHandler)
			auth.DELETE("clearCacheKey/:cacheName/:cacheKey", controllers.ClearCacheKeyHandler)
			auth.DELETE("clearCacheAll", controllers.ClearCacheAllHandler)

		}
	}
	v12 := r.Group("monitor/logininfor")
	{
		auth := v12.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("/list", controllers.LoginInformListHandler)
			auth.POST("/export", controllers.ExportHandler)
			auth.DELETE("/:infoIds", controllers.DeleteByIdHandler)
			auth.DELETE("/clean", controllers.CleanHandler)
			auth.GET("/unlock/:userName", controllers.UnlockHandler)
		}
	}
	v13 := r.Group("monitor/job") //定时任务相关
	{
		auth := v13.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("list", controllers.ListJob)
			auth.POST("export", controllers.ExportJob)
			auth.GET(":jobId", controllers.GetJobById)
			auth.POST("", controllers.SaveJob)
			auth.PUT("", controllers.UploadJob)
			auth.PUT("changeStatus", controllers.ChangeStatus)
			auth.PUT("run", controllers.RunJob)
			auth.DELETE(":jobIds", controllers.DelectJob)
		}
	}
	v14 := r.Group("monitor/jobLog")
	{
		auth := v14.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			// 需要权限
			auth.GET("list", controllers.ListJobLog)
			auth.POST("export", controllers.ExportJobLog)
			auth.GET(":configId", controllers.GetJobLog)
			auth.DELETE(":jobLogIds", controllers.DetectJobLog)
			auth.DELETE("clean", controllers.ClearJobLog)
		}
	}

	v15 := r.Group("monitor/online")
	{
		auth := v15.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("list", controllers.ListOnLine)
			auth.DELETE("/:onlineid", controllers.DetectOnLine)
		}
	}

	v16 := r.Group("monitor/operlog")
	{
		auth := v16.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("/list", controllers.ListOperlog)
			auth.DELETE("/:operId", controllers.DelectOperlog)
			auth.DELETE("/clean", controllers.ClearOperlog)
			auth.POST("/export", controllers.ExportOperlog)
		}
	}
	v17 := r.Group("live")
	{
		auth := v17.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("/get/live/stream/online/list", controllers.GetLiveStreamOnlineList)
			auth.GET("/get/live/forbid/stream/list", controllers.GetLiveForbidStreamList)
			auth.POST("/forbid/live/stream", controllers.ForbidLiveStream)
			auth.GET("/drop/live/stream", controllers.DropLiveStream)
			auth.GET("/get/live/stream/state", controllers.GetLiveStreamState)
			auth.GET("/resume/live/stream", controllers.ResumeLiveStream)
			auth.GET("/get/live/delay/info/list", controllers.GetLiveDelayInfoList)
			auth.GET("/add/delay/live/stream", controllers.AddDelayLiveStream)
			auth.GET("/resume/delay/live/stream", controllers.ResumeDelayLiveStream)
		}
	}

	v18 := r.Group("livetool")
	{
		auth := v18.Group("")
		auth.Use(middlewares.JWTAuthMiddleware())
		{
			auth.GET("/get/push/url", controllers.GetPushUrl)
			auth.GET("/get/pull/url", controllers.GetPullUrl)
		}
	}

	r.GET("/ws", controllers.Run) //这个是聊天室，目前存疑，不知道演示出来效果，并且代码需要修改
	r.POST("/articles/poster/:address", controllers.GenerateArticlePoster)
	return r
}
