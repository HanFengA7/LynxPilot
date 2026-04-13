package bootstrap

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	userrepo "github.com/LychApe/LynxPilot/internal/repository/user"
	userservice "github.com/LychApe/LynxPilot/internal/service/user"
	"github.com/LychApe/LynxPilot/internal/utils/logger"
	"golang.org/x/term"
	"gorm.io/gorm"
)

func InitUserFromTerminal(db *gorm.DB, cfg *Config) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	username = strings.TrimSpace(username)
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}

	fmt.Print("password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Println()

	fmt.Print("confirm password: ")
	confirmBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	fmt.Println()

	password := strings.TrimSpace(string(passwordBytes))
	confirm := strings.TrimSpace(string(confirmBytes))
	if password == "" {
		return fmt.Errorf("password cannot be empty")
	}
	if password != confirm {
		return fmt.Errorf("passwords do not match")
	}

	userRepository := userrepo.New(db)
	userService := userservice.New(userRepository, cfg.Auth.TokenSalt)
	if err := userService.CreateUser(username, password); err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("用户 %s 创建成功", username))
	return nil
}
