package helper

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gogf/gf/v2/frame/g"
)

// InitFunctions 初始化函数
type InitFunctions struct {
	// InitServices 初始化服务层
	InitServices func()

	// CleanupServices 清理服务资源
	CleanupServices func()
}

var (
	// 全局初始化函数
	globalInitFunctions InitFunctions
)

// SetInitFunctions 设置初始化函数
func SetInitFunctions(fns InitFunctions) {
	globalInitFunctions = fns
}

// InitAll 执行所有初始化
func InitAll() {
	if globalInitFunctions.InitServices != nil {
		globalInitFunctions.InitServices()
	}

	// 设置退出时的资源清理
	setupCleanup()
}

// setupCleanup 设置程序退出时的清理函数
func setupCleanup() {
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		g.Log().Info(ctx, "接收到退出信号，开始清理资源...")

		if globalInitFunctions.CleanupServices != nil {
			globalInitFunctions.CleanupServices()
		}

		g.Log().Info(ctx, "资源清理完成，程序退出")
		os.Exit(0)
	}()
}
