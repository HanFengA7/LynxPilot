package user

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"size:64;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:128;not null"`
}
