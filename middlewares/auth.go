package middlewares

import (
	"ruoyi/controllers"
	"ruoyi/dao/redis"
	"ruoyi/pkg/JWT"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件，token认证
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			/*c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请求头中auth为空",
			})
			*/
			controllers.ResponseError(c, controllers.CodeNeedLogin)
			c.Abort() //ctx.Abort()方法的作用 终止调用整个链条,直接退出该r.get
			return
		}
		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			/*c.JSON(http.StatusOK, gin.H{
				"code": 2004,
				"msg":  "请求头中auth格式有误",
			})*/
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := JWT.ParseToken(parts[1])
		if err != nil {
			/*c.JSON(http.StatusOK, gin.H{
				"code": 2005,
				"msg":  "无效的Token",
			})*/
			controllers.ResponseError(c, controllers.CodeInvalidToken)
			c.Abort()
			return
		}
		token, err := redis.GetLogin(mc.UserID)
		if err != nil || parts[1] != token {
			controllers.ResponseErrorWithMsg(c, controllers.CodeServerBusy, "用户已在别处登录")
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set("userID", mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get("CtxUserIDKey")来获取当前请求的用户信息
	}
}
