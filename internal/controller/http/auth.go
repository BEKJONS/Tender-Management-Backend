package http

import (
	_ "tender_management/docs"

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

	res, err := a.us.LogIn(req)
	if err != nil {
		a.log.Error("Error in login", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Success 200 {object} entity.RegisterRes
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

	res, err := a.us.AddUser(req)
	if err != nil {
		a.log.Error("Error in creating user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

//
//// UpdateUser godoc
//// @Summary Update User
//// @Description Update user details
//// @Tags User
//// @Accept json
//// @Produce json
//// @Param id path string true "User ID"
//// @Param UpdateUser body entity.UserUpdate true "Update user"
//// @Success 200 {object} entity.UserRequest
//// @Failure 400 {object} entity.Error
//// @Failure 500 {object} entity.Error
//// @Router /auth/update/{id} [put]
//func (a *authRoutes) updateUser(c *gin.Context) {
//	var req entity.UserRequest
//	var user entity.UserUpdate
//
//	if err := c.ShouldBindJSON(&user); err != nil {
//		a.log.Error("Error in getting from body", "error", err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	req.UserID = c.Param("id")
//	req.FirstName = user.FirstName
//	req.LastName = user.LastName
//	req.PhoneNumber = user.PhoneNumber
//	req.Email = user.Email
//	req.Role = user.Role
//
//	res, err := a.us.UpdateUser(req)
//	if err != nil {
//		a.log.Error("Error in updating user", "error", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, res)
//}
//
//// DeleteUser godoc
//// @Summary Delete User
//// @Description Delete a user account
//// @Tags User
//// @Accept json
//// @Produce json
//// @Param id path string true "User ID"
//// @Success 200 {object} entity.Message
//// @Failure 400 {object} entity.Error
//// @Failure 500 {object} entity.Error
//// @Router /auth/delete/{id} [delete]
//func (a *authRoutes) deleteUser(c *gin.Context) {
//	var req entity.UserID
//
//	id := c.Param("id")
//	req.ID = id
//
//	res, err := a.us.DeleteUser(req)
//
//	if err != nil {
//		a.log.Error("Error in deleting user", "error", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, res)
//}
//
//// GetUser godoc
//// @Summary Get User
//// @Description Retrieve user details by ID
//// @Tags User
//// @Accept json
//// @Produce json
//// @Param id path string true "User ID"
//// @Success 200 {object} entity.UserRequest
//// @Failure 400 {object} entity.Error
//// @Failure 500 {object} entity.Error
//// @Router /auth/get/{id} [get]
//func (a *authRoutes) getUser(c *gin.Context) {
//	var req entity.UserID
//
//	id := c.Param("id")
//	req.ID = id
//
//	res, err := a.us.GetUser(req)
//
//	if err != nil {
//		a.log.Error("Error in getting user", "error", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, res)
//}
//
//// ListUser godoc
//// @Summary List Users
//// @Description Retrieve a list of users with optional filters
//// @Tags User
//// @Accept json
//// @Produce json
//// @Param FilterUser query entity.FilterUser false "User filter parameters"
//// @Success 200 {array} entity.UserList
//// @Failure 400 {object} entity.Error
//// @Failure 500 {object} entity.Error
//// @Router /auth/list [get]
//func (a *authRoutes) listUser(c *gin.Context) {
//	var req entity.FilterUser
//
//	if err := c.ShouldBindQuery(&req); err != nil {
//		a.log.Error("Error in getting from body", "error", err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	res, err := a.us.GetUserList(req)
//
//	if err != nil {
//		a.log.Error("Error in getting user", "error", err)
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//
//	c.JSON(http.StatusOK, res)
//}
