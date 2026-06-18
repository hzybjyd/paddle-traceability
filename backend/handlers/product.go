package handlers

import (
	"net/http"
	"strconv"

	"paddle-traceability/services"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ProductHandler struct {
	productService *services.ProductService
	traceService   *services.TraceService
}

func NewProductHandler(productService *services.ProductService, traceService *services.TraceService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
		traceService:   traceService,
	}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "FACTORY" {
		c.JSON(http.StatusForbidden, Response{
			Code:    403,
			Message: "only FACTORY role can create products",
		})
		return
	}

	var req services.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid request: " + err.Error(),
		})
		return
	}

	userID := c.MustGet("user_id").(uint)

	product, txRecord, err := h.productService.CreateProduct(userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Code:    201,
		Message: "product created and recorded on chain",
		Data: map[string]interface{}{
			"product_id":   product.ID,
			"product_uid":  product.ProductUID,
			"tx_hash":      txRecord.TxHash,
			"block_height": txRecord.BlockHeight,
			"created_at":   product.CreatedAt,
		},
	})
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	userID := c.MustGet("user_id").(uint)
	role := c.MustGet("role").(string)

	products, total, err := h.productService.ListProducts(userID, role, page, pageSize)
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
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"items":     products,
		},
	})
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productService.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 200,
		Data: product,
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	role := c.MustGet("role").(string)
	if role != "FACTORY" && role != "LOGISTICS" && role != "RETAILER" {
		c.JSON(http.StatusForbidden, Response{
			Code:    403,
			Message: "permission denied",
		})
		return
	}

	idOrUID := c.Param("id")
	if idOrUID == "" {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "product id or uid is required",
		})
		return
	}

	var req services.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "invalid request",
		})
		return
	}

	userID := c.MustGet("user_id").(uint)

	product, err := h.productService.UpdateProduct(idOrUID, req, userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "product status updated and recorded on chain",
		Data:    product,
	})
}

func (h *ProductHandler) GetTrace(c *gin.Context) {
	id := c.Param("id")

	// Get product info first
	product, err := h.productService.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: err.Error(),
		})
		return
	}

	traceResult, err := h.traceService.GetFullTrace(product.ProductUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code: 200,
		Data: traceResult,
	})
}
