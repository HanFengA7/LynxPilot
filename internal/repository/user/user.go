package user

import (
	usermodel "github.com/LychApe/LynxPilot/internal/model/user"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByUsername(username string) (*usermodel.User, error) {
	var user usermodel.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) Create(user *usermodel.User) error {
	return r.db.Create(user).Error
}
