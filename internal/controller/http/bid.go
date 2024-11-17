package http

import (
	"log/slog"
	"net/http"
	"strconv"
	"tender_management/config"
	"tender_management/internal/email"
	"tender_management/internal/entity"
	"tender_management/internal/usecase"
	"tender_management/internal/usecase/repo"
	"tender_management/internal/web"
	"tender_management/pkg/postgres"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type bidRoutes struct {
	us  *usecase.BidService
	log *slog.Logger
}

func newBidRoutes(router *gin.RouterGroup, us *usecase.BidService, log *slog.Logger) {
	bids := bidRoutes{us, log}

	router.POST("/tenders/:id/bids", bids.submitBid)
	router.GET("/tenders/:id/bids", bids.getSubmittedBids)
}

// SubmitBid godoc
// @Summary Submit a bid on a tender
// @Description Contractors can submit bids on open tenders
// @Tags Bids
// @Accept json
// @Produce json
// @Param tender_id path string true "Tender ID"
// @Param bid body entity.Bid1 true "Bid details"
// @Success 201 {object} entity.Bid
// @Failure 400 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /tenders/{tender_id}/bids [post]
func (b *bidRoutes) submitBid(c *gin.Context) {
	var bid entity.Bid
	claims, err := extractClaims(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	userID, ok := claims["id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user id"})
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "cannot get user id"})
		return
	}

	tenderID := c.Param("id")

	if err := c.ShouldBindJSON(&bid); err != nil {
		b.log.Error("Error in getting from body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the tender ID from the URL parameter
	bid.TenderID = tenderID

	// Submit the bid
	newBid, err := b.us.SubmitBid(bid)
	if err != nil {
		b.log.Error("Error in submitting bid", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Send an email to the client
	message, err := email.CreateBidMessage(userID, tenderID, username)
	if err != nil {
		b.log.Error("Error in sending email", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	NotifyAll(c, message, config.NewConfig())

	// Return the submitted bid
	c.JSON(http.StatusCreated, newBid)
}

// GetSubmittedBids godoc
// @Summary Get Bids for Tender
// @Description Get a list of bids for a tender with optional filters for price, delivery time, and comments or status.
// @Tags Bids
// @Accept json
// @Produce json
// @Param price query float64 false "Filter by price"
// @Param delivery_time query int false "Filter by delivery time"
// @Param comments query string false "Filter by comments"
// @Param status query string false "Filter by status"
// @Param client_id query string false "Client ID to filter tenders"
// @Success 200 {array} entity.Bid
// @Failure 400 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /tenders/{id}/bids [get]
func (b *bidRoutes) getSubmittedBids(c *gin.Context) {
	tenderID := c.Param("id")

	req := entity.ListBidReq{
		TenderID: tenderID,
	}

	price, err := parseFloatQueryParam(c, "price")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price filter"})
		return
	}
	if price != nil {
		req.PriceFilter = price
	}

	deliveryTime, err := parseIntQueryParam(c, "delivery_time")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid delivery_time filter"})
		return
	}
	if deliveryTime != nil {
		req.DeliveryTimeFilter = deliveryTime
	}

	status := c.DefaultQuery("status", "")
	if status != "" {
		req.Status = status
	}

	comments := c.DefaultQuery("comments", "")
	if comments != "" {
		req.Comments = comments
	}

	clientID := c.DefaultQuery("client_id", "")
	if clientID != "" {
		req.ClientID = clientID
	}

	bids, err := b.us.GetBids(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bids)
}

func parseFloatQueryParam(c *gin.Context, param string) (*float64, error) {
	queryVal := c.DefaultQuery(param, "")
	if queryVal == "" {
		return nil, nil
	}
	val, err := strconv.ParseFloat(queryVal, 64)
	if err != nil {
		return nil, err
	}
	return &val, nil
}

func parseIntQueryParam(c *gin.Context, param string) (*int, error) {
	queryVal := c.DefaultQuery(param, "")
	if queryVal == "" {
		return nil, nil
	}
	val, err := strconv.Atoi(queryVal)
	if err != nil {
		return nil, err
	}
	return &val, nil
}

func NotifyAll(c *gin.Context, message string, cfg *config.Config) {
	db, err := postgres.Connection(cfg)
	if err != nil {
		return
	}
	emails, err := repo.GetAllEmails(db)
	if err != nil {
		return
	}
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	for _, email := range emails {
		web.SendNotification(c, message, cfg, client, email.Email)
	}
}
