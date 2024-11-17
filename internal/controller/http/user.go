package http

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"tender_management/internal/usecase"
)

type userController struct {
	log *slog.Logger
	ts  *usecase.TenderService
	bid *usecase.BidService
}

func newUserController(router *gin.RouterGroup, ts *usecase.TenderService, casbin *casbin.Enforcer, bid *usecase.BidService, log *slog.Logger) {
	user := userController{log, ts, bid}
	router.Use(PermissionMiddleware(casbin))
	router.GET("/:id/tenders", user.getUserTenders)
	router.GET("/:id/bids", user.getUserBids)
}

// GetUserTenders godoc
// @Summary Get User Tenders
// @Description Get all tenders associated with a specific user.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {array} entity.Tender
// @Failure 400 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /users/{id}/tenders [get]
func (u *userController) getUserTenders(c *gin.Context) {
	userID := c.MustGet("user_id").(string) // Получение ID пользователя из пути
	if userID == "" {
		u.log.Error("User ID is missing in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// Получение тендеров через сервис
	tenders, err := u.ts.GetUserTenders(userID)
	if err != nil {
		u.log.Error("Error in retrieving tenders for user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error in retrieving tenders"})
		return
	}

	c.JSON(http.StatusOK, tenders)
}

// GetUserBids godoc
// @Summary Get User Bids
// @Description Retrieve all bids placed by a specific user.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {array} entity.Bid
// @Failure 400 {object} entity.Error
// @Failure 404 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /users/{id}/bids [get]
func (u *userController) getUserBids(c *gin.Context) {
	userID := c.MustGet("user_id").(string) // Получение ID пользователя из пути
	if userID == "" {
		u.log.Error("User ID is missing in request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// Вызов метода сервиса для получения ставок
	bids, err := u.bid.GetUserBids(userID)
	if err != nil {
		u.log.Error("Error retrieving user bids", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bids)
}
