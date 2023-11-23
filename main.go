package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"ruoyi/controllers"
	"ruoyi/dao/mysql"
	"ruoyi/dao/redis"
	"ruoyi/logger"
	"ruoyi/pkg/snowflake"
	"ruoyi/routes"
	"ruoyi/scheduler"
	"ruoyi/settings"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func main() {
	// 1. 加载配置
	if err := settings.Init(); err != nil {
		fmt.Println("init Settings failed!!! err:", err)
		return
	}
	// 2. 初始化日志
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.AppConfig.Mode); err != nil {
		fmt.Println("init Logger failed!!! err : ", err)
		return
	}
	defer zap.L().Sync() // 将缓存区的日志追加到磁盘文件当中
	zap.L().Debug("logger init success.....")
	// 3. 初始化SQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Println("init Mysql failed!!! err : ", err)
		return
	}
	defer mysql.Close()
	// 4. 初始化redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Println("init Redis failed!!! err", err)
		return
	}
	defer redis.Close()

	scheduler.InitCron()

	if err := snowflake.Init(settings.Conf.AppConfig.StartTime, settings.Conf.AppConfig.MachineID); err != nil {
		fmt.Printf("init snowflake failed,err:%v\n", err)
		return
	}
	//初始化gin框架内置的校验器使用的翻译器
	if err := controllers.InitTrans("zh"); err != nil {
		fmt.Printf("init validator Trans failed,err:%v\n", err)
	}
	// 5. 注册路由
	r := routes.SetUp(settings.Conf.AppConfig.Mode, settings.Conf.RedisConfig.Host, settings.Conf.RedisConfig.Port)

	// 6.启动服务（优雅关机，重启）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.AppConfig.Port),
		Handler: r,
	}

	go func() {
		// 开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: %s\n", zap.Error(err))
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1) // 创建一个接收信号的通道
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit                                               // 阻塞在此，当接收到上述两种信号时才会往下执行
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	mysql.DeleteOnlineAll()
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}

	zap.L().Info("Server exiting")
}
