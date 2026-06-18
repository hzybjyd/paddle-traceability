package handlers

import (
	"net/http"

	"paddle-traceability/services"

	"github.com/gin-gonic/gin"
)

type VerifyHandler struct {
	traceService *services.TraceService
}

func NewVerifyHandler(traceService *services.TraceService) *VerifyHandler {
	return &VerifyHandler{traceService: traceService}
}

func (h *VerifyHandler) VerifyProduct(c *gin.Context) {
	productUID := c.Param("product_uid")
	if productUID == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "missing product_uid",
		})
		return
	}

	result, err := h.traceService.VerifyProduct(productUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "error occurred during verification",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":     200,
		"verified": result.Verified,
		"message":  result.Message,
		"data":     result.Data,
	})
}
