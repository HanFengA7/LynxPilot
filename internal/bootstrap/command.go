package bootstrap

import (
	"fmt"

	"gorm.io/gorm"
)

func HandleCommand(args []string, db *gorm.DB, cfg *Config) (bool, error) {
	if len(args) < 2 {
		return false, nil
	}

	switch args[1] {
	case "init-user":
		if err := InitUserFromTerminal(db, cfg); err != nil {
			return true, err
		}
		return true, nil
	default:
		return true, fmt.Errorf("unknown command: %s", args[1])
	}
}
