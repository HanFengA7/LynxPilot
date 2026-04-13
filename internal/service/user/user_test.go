package user

import (
	"testing"

	usermodel "github.com/LychApe/LynxPilot/internal/model/user"
	userrepo "github.com/LychApe/LynxPilot/internal/repository/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestLoginTrimsPasswordLikeCreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open test db: %v", err)
	}

	if err := db.AutoMigrate(&usermodel.User{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}

	repo := userrepo.New(db)
	service := New(repo, "test-salt")

	if err := service.CreateUser("admin", " secret "); err != nil {
		t.Fatalf("create user: %v", err)
	}

	token, err := service.Login("admin", " secret ")
	if err != nil {
		t.Fatalf("login with spaced password: %v", err)
	}

	if token == "" {
		t.Fatal("expected token")
	}
}
