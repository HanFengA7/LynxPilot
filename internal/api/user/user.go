package user

import (
	"net/http"

	userservice "github.com/LychApe/LynxPilot/internal/service/user"
	"github.com/gin-gonic/gin"
)

// API 用户API
type API struct {
	userService *userservice.Service
}

// loginRequest 登录请求结构体
type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// New 创建用户API
func New(userService *userservice.Service) *API {
	return &API{userService: userService}
}

// Login 登录接口
func (a *API) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	token, err := a.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
