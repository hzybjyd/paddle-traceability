package handlers

import (
	"net/http"
	"time"

	"paddle-traceability/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterRequest struct {
	Username    string `json:"username" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Role        string `json:"role" binding:"required"`
	CompanyName string `json:"company_name"`
	Phone       string `json:"phone"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid request: " + err.Error(),
		})
		return
	}

	validRoles := map[string]bool{
		"FACTORY": true, "LOGISTICS": true, "RETAILER": true,
	}
	if !validRoles[req.Role] {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid role type",
		})
		return
	}

	user, err := h.authService.Register(req.Username, req.Password, req.Role, req.CompanyName, req.Phone)
	if err != nil {
		c.JSON(http.StatusConflict, Response{
			Code:    409,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "register success",
		Data: map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid request",
		})
		return
	}

	token, user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, Response{
			Code:    401,
			Message: "username or password incorrect",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "login success",
		Data: map[string]interface{}{
			"token":      token,
			"expires_at": time.Now().Add(24 * time.Hour).Format(time.RFC3339),
			"user": map[string]interface{}{
				"user_id":      user.ID,
				"username":     user.Username,
				"role":         user.Role,
				"company_name": user.CompanyName,
			},
		},
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 200,
		Data: map[string]interface{}{
			"user_id":      user.ID,
			"username":     user.Username,
			"role":         user.Role,
			"company_name": user.CompanyName,
			"phone":        user.Phone,
			"created_at":   user.CreatedAt,
		},
	})
}
