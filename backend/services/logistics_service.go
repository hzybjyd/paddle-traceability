package services

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"paddle-traceability/blockchain"
	"paddle-traceability/database"
	"paddle-traceability/models"

	"gorm.io/gorm"
)

type LogisticsService struct {
	chain *blockchain.ChainClient
}

func NewLogisticsService(chain *blockchain.ChainClient) *LogisticsService {
	return &LogisticsService{chain: chain}
}

type AddLogisticsRequest struct {
	ProductUID    string `json:"product_uid" binding:"required"`
	Action        string `json:"action" binding:"required"`
	WarehouseName string `json:"warehouse_name" binding:"required"`
	Location      string `json:"location" binding:"required"`
	Carrier       string `json:"carrier"`
	Remark        string `json:"remark"`
}

func (s *LogisticsService) AddRecord(operatorID uint, operatorRole string, req AddLogisticsRequest) (*models.LogisticsRecord, error) {
	var product models.Product
	if err := database.DB.Where("product_uid = ?", req.ProductUID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, errors.New("query product failed")
	}

	record := &models.LogisticsRecord{
		ProductID:     product.ID,
		Action:        req.Action,
		OperatorID:    operatorID,
		WarehouseName: req.WarehouseName,
		Location:      req.Location,
		Carrier:       req.Carrier,
		Remark:        req.Remark,
		CreatedAt:     time.Now(),
	}

	// Compute data hash
	dataBytes, _ := json.Marshal(map[string]interface{}{
		"product_uid":    req.ProductUID,
		"action":         req.Action,
		"warehouse_name": req.WarehouseName,
		"location":       req.Location,
	})
	dataHash := fmt.Sprintf("%x", sha256.Sum256(dataBytes))

	// Update product status
	newStatus := "IN_TRANSIT"
	if req.Action == "INBOUND" {
		newStatus = "IN_STOCK"
	}

	// Build updateGoods reason JSON
	reason := map[string]string{
		"product_uid":    req.ProductUID,
		"action":         req.Action,
		"new_status":     newStatus,
		"data_hash":      dataHash,
		"operator_role":  operatorRole,
		"warehouse_name": req.WarehouseName,
		"location":       req.Location,
	}
	reasonJSON, _ := json.Marshal(reason)

	// Submit to chain
	chainResult, err := s.chain.UpdateGoods(req.ProductUID, string(reasonJSON))
	if err != nil {
		return nil, errors.New("chain attestation failed")
	}

	tx := database.DB.Begin()

	if err := tx.Create(record).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("create logistics record failed")
	}

	product.Status = newStatus
	product.UpdatedAt = time.Now()
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("update product status failed")
	}

	txRecord := &models.TxRecord{
		ProductID:   product.ID,
		TxHash:      chainResult.TxHash,
		TxType:      "TRANSFER",
		DataHash:    dataHash,
		ChainStatus: "CONFIRMED",
		BlockHeight: chainResult.BlockHeight,
		OperatorID:  operatorID,
		CreatedAt:   time.Now(),
	}
	if err := tx.Create(txRecord).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("save tx record failed")
	}

	tx.Commit()
	return record, nil
}

func (s *LogisticsService) GetRecords(productUID string) ([]models.LogisticsRecord, error) {
	var product models.Product
	if err := database.DB.Where("product_uid = ?", productUID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, errors.New("query product failed")
	}

	var records []models.LogisticsRecord
	if err := database.DB.Where("product_id = ?", product.ID).Order("created_at ASC").Find(&records).Error; err != nil {
		return nil, errors.New("query logistics records failed")
	}

	return records, nil
}
