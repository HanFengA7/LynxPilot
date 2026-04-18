package server

import (
	"os"
	"syscall"
	"time"

	processService "github.com/LychApe/LynxPilot/internal/service/process"
	"github.com/LychApe/LynxPilot/internal/utils/logger"
	"github.com/gin-gonic/gin"
)

func RebootHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "reboot signal accepted"})

	go func() {
		// 给响应一点时间写回客户端，再触发进程退出
		time.Sleep(200 * time.Millisecond)

		if err := processService.StartNewProcess(); err != nil {
			logger.Errorf("拉起新进程失败: %v", err)
			return
		}

		pid := os.Getpid()
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			logger.Errorf("发送重启信号失败: %v", err)
			os.Exit(1)
		}
	}()
}

func ShutdownHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "shutdown signal accepted"})

	go func() {
		pid := os.Getpid()
		if err := syscall.Kill(pid, syscall.SIGTERM); err != nil {
			logger.Errorf("发送关闭信号失败: %v", err)
			os.Exit(1)
		}
	}()
}
