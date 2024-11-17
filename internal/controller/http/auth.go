package http

import (
	_ "tender_management/docs"

	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"tender_management/internal/entity"
	"tender_management/internal/usecase"
)

type authRoutes struct {
	us  *usecase.UserUseCase
	log *slog.Logger
}

func newUserRoutes(router *gin.RouterGroup, us *usecase.UserUseCase, log *slog.Logger) {

	auth := authRoutes{us, log}

	router.POST("/user/register", auth.createUser)
	router.POST("/login", auth.login)

}

// ------------ Handler methods --------------------------------------------------------

// Login godoc
// @Summary Admin Login
// @Description Login for admin users
// @Tags User
// @Accept json
// @Produce json
// @Param Login body entity.LogInReq true "Admin login"
// @Success 200 {object} entity.LogInRes
// @Failure 400 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /auth/login [post]
func (a *authRoutes) login(c *gin.Context) {
	var req entity.LogInReq

	if err := c.ShouldBindJSON(&req); err != nil {
		a.log.Error("Error in getting from body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Username == "" || req.Password == "" {
		a.log.Error("Unauthorized login attempt", "error", "Username and password are required")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	res, err := a.us.LogIn(req)
	if err != nil {
		switch err.Error() {
		case "Invalid username or password":
			a.log.Error("Unauthorized login attempt", "error", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		case "failed to get user: sql: no rows in result set":
			a.log.Error("User not found", "error", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		default:
			a.log.Error("Error in login", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, res)
}

// CreateUser godoc
// @Summary Create User
// @Description Register a new user account
// @Tags User
// @Accept json
// @Produce json

// @Param CreateUser body entity.RegisterReq true "Create user"
// @Success 200 {object} entity.LogInRes
// @Failure 400 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /auth/user/register [post]
func (a *authRoutes) createUser(c *gin.Context) {
	var req entity.RegisterReq

	if err := c.ShouldBindJSON(&req); err != nil {
		a.log.Error("Error in getting from body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := checkmail.ValidateFormat(req.Email)
	if err != nil {
		a.log.Error("Invalid email provided", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		return
	}
	if req.Username == "" || req.Email == "" {
		a.log.Error("message", "username or email cannot be empty")
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or email cannot be empty"})
		return
	}
	if req.Role != "contractor" && req.Role != "client" {
		a.log.Error("Invalid role provided", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}
	exists, err := a.us.IsEmailExists(req.Email)
	if err != nil {
		a.log.Error("Error in checking if email exists", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if exists {
		a.log.Error("Email already exists", "email", req.Email)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already exists"})
		return
	}

	res, err := a.us.AddUser(req)
	if err != nil {
		a.log.Error("Error in creating user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, res)
}
