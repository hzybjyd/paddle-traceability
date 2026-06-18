package handlers

import (
	"net/http"

	"paddle-traceability/services"

	"github.com/gin-gonic/gin"
)

type LogisticsHandler struct {
	logisticsService *services.LogisticsService
}

func NewLogisticsHandler(logisticsService *services.LogisticsService) *LogisticsHandler {
	return &LogisticsHandler{logisticsService: logisticsService}
}

func (h *LogisticsHandler) AddRecord(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "LOGISTICS" && role != "RETAILER" {
		c.JSON(http.StatusForbidden, Response{
			Code:    403,
			Message: "only LOGISTICS and RETAILER can add logistics records",
		})
		return
	}

	var req services.AddLogisticsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid request: " + err.Error(),
		})
		return
	}

	userID := c.MustGet("user_id").(uint)

	record, err := h.logisticsService.AddRecord(userID, role, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "logistics record added and recorded on chain",
		Data: map[string]interface{}{
			"logistics_id": record.ID,
			"product_uid":  req.ProductUID,
			"action":       record.Action,
			"created_at":   record.CreatedAt,
		},
	})
}

func (h *LogisticsHandler) GetRecords(c *gin.Context) {
	productUID := c.Query("product_uid")
	if productUID == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "missing product_uid parameter",
		})
		return
	}

	records, err := h.logisticsService.GetRecords(productUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 200,
		Data: map[string]interface{}{
			"product_uid": productUID,
			"records":     records,
		},
	})
}
