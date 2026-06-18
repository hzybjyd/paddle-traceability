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
	"paddle-traceability/utils"

	"gorm.io/gorm"
)

type ProductService struct {
	chain *blockchain.ChainClient
}

func NewProductService(chain *blockchain.ChainClient) *ProductService {
	return &ProductService{chain: chain}
}

type CreateProductRequest struct {
	Brand             string `json:"brand" binding:"required"`
	Model             string `json:"model" binding:"required"`
	Material          string `json:"material" binding:"required"`
	RubberType        string `json:"rubber_type"`
	BatchNo           string `json:"batch_no" binding:"required"`
	ProductionDate    string `json:"production_date" binding:"required"`
	QualityReportHash string `json:"quality_report_hash"`
}

type UpdateProductRequest struct {
	Status string `json:"status"`
	Remark string `json:"remark"`
}

func (s *ProductService) CreateProduct(factoryID uint, req CreateProductRequest) (*models.Product, *models.TxRecord, error) {
	prodDate, err := time.Parse("2006-01-02", req.ProductionDate)
	if err != nil {
		return nil, nil, errors.New("invalid production date format, expected YYYY-MM-DD")
	}

	productUID := utils.GenerateID()

	product := &models.Product{
		ProductUID:        productUID,
		Brand:             req.Brand,
		Model:             req.Model,
		Material:          req.Material,
		RubberType:        req.RubberType,
		BatchNo:           req.BatchNo,
		ProductionDate:    prodDate,
		QualityReportHash: req.QualityReportHash,
		FactoryID:         factoryID,
		Status:            "PRODUCED",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Compute product data hash
	dataBytes, _ := json.Marshal(map[string]interface{}{
		"product_uid": productUID,
		"brand":       req.Brand,
		"model":       req.Model,
		"material":    req.Material,
		"batch_no":    req.BatchNo,
		"prod_date":   req.ProductionDate,
	})
	dataHash := fmt.Sprintf("%x", sha256.Sum256(dataBytes))

	// Build createGoods desc JSON
	desc := map[string]string{
		"product_uid":   productUID,
		"data_hash":     dataHash,
		"operator_role": "FACTORY",
	}
	descJSON, _ := json.Marshal(desc)

	// Invoke blockchain
	chainResult, err := s.chain.CreateGoods(productUID, string(descJSON))
	if err != nil {
		return nil, nil, errors.New("chain attestation failed")
	}

	// Transaction: save product and tx record
	tx := database.DB.Begin()
	if err := tx.Create(product).Error; err != nil {
		tx.Rollback()
		return nil, nil, errors.New("create product failed")
	}

	txRecord := &models.TxRecord{
		ProductID:   product.ID,
		TxHash:      chainResult.TxHash,
		TxType:      "CREATE",
		DataHash:    dataHash,
		ChainStatus: "CONFIRMED",
		BlockHeight: chainResult.BlockHeight,
		OperatorID:  factoryID,
		CreatedAt:   time.Now(),
	}
	if err := tx.Create(txRecord).Error; err != nil {
		tx.Rollback()
		return nil, nil, errors.New("save tx record failed")
	}

	tx.Commit()
	return product, txRecord, nil
}

func (s *ProductService) GetProduct(idOrUID string) (*models.Product, error) {
	var product models.Product
	err := database.DB.Where("id = ? OR product_uid = ?", idOrUID, idOrUID).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, errors.New("query product failed")
	}
	return &product, nil
}

func (s *ProductService) ListProducts(userID uint, role string, page, pageSize int) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := database.DB.Model(&models.Product{})
	// Only FACTORY users can view the products they produced.
	// Other authenticated roles (LOGISTICS/RETAILER) can view all products for operation needs.
	if role == "FACTORY" {
		query = query.Where("factory_id = ?", userID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.New("query product list failed")
	}

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&products).Error
	if err != nil {
		return nil, 0, errors.New("query product list failed")
	}

	return products, total, nil
}

func (s *ProductService) UpdateProduct(idOrUID string, req UpdateProductRequest, operatorID uint, operatorRole string) (*models.Product, error) {
	var product models.Product
	if err := database.DB.Where("id = ? OR product_uid = ?", idOrUID, idOrUID).First(&product).Error; err != nil {
		return nil, errors.New("product not found")
	}

	if req.Status != "" {
		product.Status = req.Status
	}
	product.UpdatedAt = time.Now()

	// Compute data hash
	dataBytes, _ := json.Marshal(map[string]interface{}{
		"product_uid": product.ProductUID,
		"status":      product.Status,
	})
	dataHash := fmt.Sprintf("%x", sha256.Sum256(dataBytes))

	// Build updateGoods reason JSON
	reason := map[string]string{
		"new_status":    product.Status,
		"data_hash":     dataHash,
		"operator_role": operatorRole,
	}
	reasonJSON, _ := json.Marshal(reason)

	// Submit to chain
	chainResult, err := s.chain.UpdateGoods(product.ProductUID, string(reasonJSON))
	if err != nil {
		return nil, errors.New("chain attestation failed")
	}

	tx := database.DB.Begin()
	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("update product failed")
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
	if req.Status == "SOLD" {
		txRecord.TxType = "CONFIRM"
	}
	if err := tx.Create(txRecord).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("save tx record failed")
	}

	tx.Commit()
	return &product, nil
}
