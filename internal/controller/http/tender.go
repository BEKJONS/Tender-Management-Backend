package http

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	_ "tender_management/docs"
	"tender_management/internal/entity"
	"tender_management/internal/usecase"
)

type tenderRoutes struct {
	ts  *usecase.TenderService
	log *slog.Logger
}

func newTenderRoutes(router *gin.RouterGroup, ts *usecase.TenderService, casbin *casbin.Enforcer, log *slog.Logger) {

	tender := tenderRoutes{ts, log}
	router.Use(PermissionMiddleware(casbin))
	router.POST("/", tender.createTender)
	router.GET("/", tender.listTenders)
	router.PUT("/:id/:status", tender.updateTenderStatus)
	router.DELETE("/:id", tender.deleteTender)

	router.POST("/:id/award/:bid_id", tender.awardTender)
}

// ------------ Handler methods --------------------------------------------------------

// CreateTender godoc
// @Summary Create Tender
// @Description Create a new tender for a client with details like title, description, deadline, and budget.
// @Tags Tender
// @Accept json
// @Produce json
// @Param CreateTender body entity.TenderReq1 true "Create tender"
// @Success 201 {object} entity.Tender
// @Failure 400 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /tenders [post]
func (t *tenderRoutes) createTender(c *gin.Context) {
	var req entity.TenderReq1

	if err := c.ShouldBindJSON(&req); err != nil {
		t.log.Error("Error in getting from body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create tender via service
	tender, err := t.ts.CreateTender(entity.TenderReq{ClientID: c.MustGet("user_id").(string),
		Title: req.Title, Description: req.Description,
		Deadline: req.Deadline, Budget: req.Budget})
	if err != nil {
		t.log.Error("Error in creating tender", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, tender)
}

// ListTenders godoc
// @Summary List All Tenders
// @Description List all tenders for a specific client
// @Tags Tender
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param client_id query string false "Client ID to filter tenders"
// @Success 200 {array} entity.Tender
// @Failure 500 {object} entity.Error
// @Router /tenders [get]
func (t *tenderRoutes) listTenders(c *gin.Context) {
	clientID := c.DefaultQuery("client_id", c.MustGet("user_id").(string))

	tenders, err := t.ts.ListTenders(clientID)
	if err != nil {
		t.log.Error("Error in listing tenders", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tenders)
}

// UpdateTenderStatus godoc
// @Summary Update Tender Status
// @Description Update the status of a tender (open, closed, awarded)
// @Tags Tender
// @Accept json
// @Produce json
// @Param id path string true "Tender ID"
// @Param status path string true "Update status"
// @Success 200 {object} entity.Message
// @Failure 400 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /tenders/{id}/{status} [put]
func (t *tenderRoutes) updateTenderStatus(c *gin.Context) {
	req := entity.UpdateTender{}

	req.Id = c.Param("id")
	req.Status = c.Param("status")

	// Update status via service
	msg, err := t.ts.UpdateTenderStatus(&req)
	if err != nil {
		t.log.Error("Error in updating tender status", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, msg)
}

// DeleteTender godoc
// @Summary Delete Tender
// @Description Delete a tender by its ID
// @Tags Tender
// @Accept json
// @Produce json
// @Param id path string true "Tender ID"
// @Success 200 {object} entity.Message
// @Failure 500 {object} entity.Error
// @Router /tenders/{id} [delete]
func (t *tenderRoutes) deleteTender(c *gin.Context) {
	tenderID := c.Param("id")

	// Delete tender via service
	msg, err := t.ts.DeleteTender(tenderID)
	if err != nil {
		t.log.Error("Error in deleting tender", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, msg)
}

// awardTender godoc
// @Summary Award Tender
// @Description Award a bid to a specific tender by tender ID and bid ID.
// @Tags Tender
// @Accept json
// @Produce json
// @Param id path string true "Tender ID"
// @Param bid_id path string true "Bid ID"
// @Success 200 {object} entity.AwardedRes
// @Failure 400 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /tenders/{id}/award/{bid_id} [post]
func (t *tenderRoutes) awardTender(c *gin.Context) {
	tenderID := c.Param("id")
	bidID := c.Param("bid_id")

	if tenderID == "" || bidID == "" {
		t.log.Error("Missing path parameters", "tenderID", tenderID, "bidID", bidID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tender ID and Bid ID are required"})
		return
	}

	req := entity.Awarded{
		TenderID: tenderID,
		BideId:   bidID,
	}

	res, err := t.ts.AwardTender(&req)
	if err != nil {
		t.log.Error("Error in awarding tender", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
