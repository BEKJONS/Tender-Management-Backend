package http

import (
	_ "tender_management/docs"

	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"tender_management/internal/entity"
	"tender_management/internal/usecase"
)

type tenderRoutes struct {
	ts  *usecase.TenderService
	log *slog.Logger
}

func newTenderRoutes(router *gin.RouterGroup, ts *usecase.TenderService, log *slog.Logger) {

	tender := tenderRoutes{ts, log}

	router.POST("/tenders", tender.createTender)
	router.GET("/tenders", tender.listTenders)
	router.PUT("/tenders/:id", tender.updateTenderStatus)
	router.DELETE("/tenders/:id", tender.deleteTender)
}

// ------------ Handler methods --------------------------------------------------------

// CreateTender godoc
// @Summary Create Tender
// @Description Create a new tender for a client with details like title, description, deadline, and budget.
// @Tags Tender
// @Accept json
// @Produce json
// @Param CreateTender body entity.TenderReq true "Create tender"
// @Success 201 {object} entity.Tender
// @Failure 400 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /tenders [post]
func (t *tenderRoutes) createTender(c *gin.Context) {
	var req entity.TenderReq

	if err := c.ShouldBindJSON(&req); err != nil {
		t.log.Error("Error in getting from body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create tender via service
	tender, err := t.ts.CreateTender(req)
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
// @Param client_id query string false "Client ID to filter tenders"
// @Success 200 {array} entity.Tender
// @Failure 500 {object} entity.Error
// @Router /tenders [get]
func (t *tenderRoutes) listTenders(c *gin.Context) {
	clientID := c.DefaultQuery("client_id", "")

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
// @Param status body string true "New tender status (open, closed, awarded)"
// @Success 200 {object} entity.Message
// @Failure 400 {object} entity.Error
// @Failure 500 {object} entity.Error
// @Router /tenders/{id} [put]
func (t *tenderRoutes) updateTenderStatus(c *gin.Context) {
	tenderID := c.Param("id")
	var status entity.StatusRequest

	if err := c.ShouldBindJSON(&status); err != nil {
		t.log.Error("Error in getting from body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update status via service
	msg, err := t.ts.UpdateTenderStatus(tenderID, status.Status)
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
