package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/senoagung27/warehousex/internal/dto"
	"github.com/senoagung27/warehousex/internal/service"
)

type AuthController struct {
	authService service.AuthServiceInterface
}

func NewAuthController(authService service.AuthServiceInterface) *AuthController {
	return &AuthController{authService: authService}
}

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body dto.RegisterInput true "Register Input"
// @Success 201 {object} dto.AuthResponse
// @Router /api/v1/auth/register [post]
func (ctrl *AuthController) Register(c *gin.Context) {
	var input dto.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.authService.Register(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user registered successfully",
		"data":    response,
	})
}

// Login godoc
// @Summary Login user
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body dto.LoginInput true "Login Input"
// @Success 200 {object} dto.AuthResponse
// @Router /api/v1/auth/login [post]
func (ctrl *AuthController) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.authService.Login(input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"data":    response,
	})
}
