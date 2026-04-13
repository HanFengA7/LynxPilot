package bootstrap

import (
	"fmt"

	"github.com/LychApe/LynxPilot/internal/utils/logger"
)

func Run(args []string) error {
	configPath := ResolveConfigPath()
	cfg, err := LoadConfig(configPath)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to load config file: %w", err)
		logger.Error(wrappedErr.Error())
		return wrappedErr
	}
	logger.Info("配置加载成功")

	dbPath := ResolveSQLitePath()
	db, err := NewGorm(dbPath)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to initialize gorm: %w", err)
		logger.Error(wrappedErr.Error())
		return wrappedErr
	}
	logger.Info("数据库初始化成功")

	handled, err := HandleCommand(args, db, cfg)
	if err != nil {
		wrappedErr := fmt.Errorf("failed to handle command: %w", err)
		logger.Error(wrappedErr.Error())
		return wrappedErr
	}
	if handled {
		logger.Info("命令执行完成")
		return nil
	}

	r := NewGinEngine(db, cfg)
	logger.Info("服务启动中，监听地址 " + cfg.HTTPAddr())
	if err := r.Run(cfg.HTTPAddr()); err != nil {
		wrappedErr := fmt.Errorf("failed to start server: %w", err)
		logger.Error(wrappedErr.Error())
		return wrappedErr
	}

	return nil
}
