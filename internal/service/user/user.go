package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	usermodel "github.com/LychApe/LynxPilot/internal/model/user"
	userrepo "github.com/LychApe/LynxPilot/internal/repository/user"
	"gorm.io/gorm"
)

type Service struct {
	repo      *userrepo.Repository
	tokenSalt string
}

func New(repo *userrepo.Repository, tokenSalt string) *Service {
	return &Service{repo: repo, tokenSalt: tokenSalt}
}

func (s *Service) CreateUser(username, password string) error {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	if username == "" || password == "" {
		return fmt.Errorf("username and password cannot be empty")
	}

	_, err := s.repo.FindByUsername(username)
	if err == nil {
		return fmt.Errorf("user already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.repo.Create(&usermodel.User{
		Username:     username,
		PasswordHash: s.hashPassword(password),
	})
}

func (s *Service) Login(username, password string) (string, error) {
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)

	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", fmt.Errorf("invalid username or password")
	}

	if user.PasswordHash != s.hashPassword(password) {
		return "", fmt.Errorf("invalid username or password")
	}

	return s.generateToken(user.Username), nil
}

func (s *Service) ValidateToken(token string) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(token))
	if err != nil {
		return "", fmt.Errorf("invalid token")
	}

	parts := strings.SplitN(string(raw), ".", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid token")
	}

	payload := parts[0]
	signature := parts[1]
	expectedSignature := s.sign(payload)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return "", fmt.Errorf("invalid token")
	}

	payloadParts := strings.SplitN(payload, ":", 2)
	if len(payloadParts) != 2 || strings.TrimSpace(payloadParts[0]) == "" {
		return "", fmt.Errorf("invalid token")
	}

	// 先保留最小骨架，后续可在这里补 token 过期时间校验。
	return payloadParts[0], nil
}

func (s *Service) hashPassword(password string) string {
	sum := sha256.Sum256([]byte(s.tokenSalt + ":" + password))
	return hex.EncodeToString(sum[:])
}

func (s *Service) generateToken(username string) string {
	ts := time.Now().Unix()
	payload := fmt.Sprintf("%s:%d", username, ts)
	sig := s.sign(payload)
	raw := payload + "." + sig
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func (s *Service) sign(payload string) string {
	mac := hmac.New(sha256.New, []byte(s.tokenSalt))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}
