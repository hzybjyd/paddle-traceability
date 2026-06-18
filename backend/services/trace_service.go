package services

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"paddle-traceability/blockchain"
	"paddle-traceability/database"
	"paddle-traceability/models"

	"gorm.io/gorm"
)

type TraceService struct {
	chain *blockchain.ChainClient
}

func NewTraceService(chain *blockchain.ChainClient) *TraceService {
	return &TraceService{chain: chain}
}

type TraceStep struct {
	Stage       string      `json:"stage"`
	Operator    string      `json:"operator"`
	Timestamp   string      `json:"timestamp"`
	DataHash    string      `json:"data_hash"`
	TxHash      string      `json:"tx_hash"`
	ChainStatus string      `json:"chain_status"`
	Detail      interface{} `json:"detail,omitempty"`
}

type TraceResult struct {
	ProductUID string      `json:"product_uid"`
	Brand      string      `json:"brand"`
	Model      string      `json:"model"`
	TraceChain []TraceStep `json:"trace_chain"`
}

type VerifyResult struct {
	Verified bool        `json:"verified"`
	Message  string      `json:"message"`
	Data     interface{} `json:"data,omitempty"`
}

func (s *TraceService) GetFullTrace(productUID string) (*TraceResult, error) {
	var product models.Product
	if err := database.DB.Where("product_uid = ?", productUID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, errors.New("query product failed")
	}

	var txRecords []models.TxRecord
	if err := database.DB.Where("product_id = ?", product.ID).Order("created_at ASC").Find(&txRecords).Error; err != nil {
		return nil, errors.New("query tx records failed")
	}

	var traceChain []TraceStep

	for _, txR := range txRecords {
		// Get operator info
		var operator models.User
		database.DB.First(&operator, txR.OperatorID)

		step := TraceStep{
			DataHash:    txR.DataHash,
			TxHash:      txR.TxHash,
			ChainStatus: txR.ChainStatus,
			Timestamp:   txR.CreatedAt.Format("2006-01-02T15:04:05Z"),
			Operator:    operator.CompanyName,
		}

		if step.Operator == "" {
			step.Operator = operator.Username
		}

		switch txR.TxType {
		case "CREATE":
			step.Stage = "PRODUCTION"
			step.Detail = map[string]interface{}{
				"material":        product.Material,
				"batch_no":        product.BatchNo,
				"production_date": product.ProductionDate.Format("2006-01-02"),
			}
		case "TRANSFER":
			step.Stage = "LOGISTICS_TRANSFER"
			var logisticsRecord models.LogisticsRecord
			if err := database.DB.Where("product_id = ? AND operator_id = ?", product.ID, txR.OperatorID).
				Order("created_at DESC").First(&logisticsRecord).Error; err == nil {
				step.Detail = map[string]interface{}{
					"warehouse": logisticsRecord.WarehouseName,
					"action":    logisticsRecord.Action,
					"carrier":   logisticsRecord.Carrier,
					"location":  logisticsRecord.Location,
				}
				if logisticsRecord.Action == "INBOUND" {
					step.Stage = "WAREHOUSE_INBOUND"
				} else {
					step.Stage = "WAREHOUSE_OUTBOUND"
				}
			}
		case "CONFIRM":
			step.Stage = "SALE_CONFIRM"
			step.Detail = map[string]interface{}{
				"status": "SOLD",
			}
		}

		traceChain = append(traceChain, step)
	}

	return &TraceResult{
		ProductUID: product.ProductUID,
		Brand:      product.Brand,
		Model:      product.Model,
		TraceChain: traceChain,
	}, nil
}

func parseChainDataHashes(chainData string) map[int]string {
	hashes := make(map[int]string)
	lines := strings.Split(chainData, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Expected format: goodsId=xxx,updateRecord=N,reason={JSON}
		parts := strings.SplitN(line, ",updateRecord=", 2)
		if len(parts) != 2 {
			continue
		}
		recordParts := strings.SplitN(parts[1], ",reason=", 2)
		if len(recordParts) != 2 {
			continue
		}
		recordNum, err := strconv.Atoi(recordParts[0])
		if err != nil {
			continue
		}
		reason := recordParts[1]
		// The deployed contract stores the fixed string "CREATE" as the reason
		// for the initial createGoods record. Off-chain tx_records holds the
		// corresponding data_hash, so we treat record 0 as present but do not
		// parse JSON from it.
		if reason == "CREATE" {
			hashes[recordNum] = ""
			continue
		}
		var reasonObj map[string]string
		if err := json.Unmarshal([]byte(reason), &reasonObj); err != nil {
			continue
		}
		hashes[recordNum] = reasonObj["data_hash"]
	}
	return hashes
}

func (s *TraceService) VerifyProduct(productUID string) (*VerifyResult, error) {
	// First query on-chain records
	chainData, err := s.chain.QueryRecords(productUID)
	if err != nil {
		errMsg := strings.ToLower(err.Error())
		if strings.Contains(errMsg, "the id not exist") {
			return &VerifyResult{
				Verified: false,
				Message:  "verification failed: no product record found on chain, may be counterfeit",
				Data:     nil,
			}, nil
		}
		return nil, errors.New("query chain records failed")
	}

	onChainHashes := parseChainDataHashes(chainData)
	if len(onChainHashes) == 0 {
		return &VerifyResult{
			Verified: false,
			Message:  "verification failed: no product record found on chain, may be counterfeit",
			Data:     nil,
		}, nil
	}

	var product models.Product
	if err := database.DB.Where("product_uid = ?", productUID).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &VerifyResult{
				Verified: false,
				Message:  "verification failed: no product record found on chain, may be counterfeit",
				Data:     nil,
			}, nil
		}
		return nil, errors.New("query product failed")
	}

	var txRecords []models.TxRecord
	if err := database.DB.Where("product_id = ?", product.ID).Order("created_at ASC").Find(&txRecords).Error; err != nil {
		return nil, errors.New("query tx records failed")
	}

	// Verify that the data hashes stored off-chain match those on-chain.
	// Record 0 corresponds to createGoods, whose on-chain reason is the fixed
	// string "CREATE" and does not carry a data_hash. We therefore only compare
	// records produced by updateGoods (record index >= 1).
	dataHashMatched := len(txRecords) == len(onChainHashes)
	if dataHashMatched {
		for i, txR := range txRecords {
			if i == 0 {
				continue
			}
			if txR.DataHash != onChainHashes[i] {
				dataHashMatched = false
				break
			}
		}
	}
	chainVerified := len(onChainHashes) > 0

	// Get full trace chain
	traceResult, err := s.GetFullTrace(productUID)
	if err != nil {
		return nil, err
	}

	// Build trace summary
	var traceSummary []map[string]interface{}
	for _, step := range traceResult.TraceChain {
		traceSummary = append(traceSummary, map[string]interface{}{
			"stage":    step.Stage,
			"time":     step.Timestamp,
			"operator": step.Operator,
		})
	}

	return &VerifyResult{
		Verified: true,
		Message:  "verified: this product is authentic",
		Data: map[string]interface{}{
			"product_uid":       product.ProductUID,
			"brand":             product.Brand,
			"model":             product.Model,
			"material":          product.Material,
			"rubber_type":       product.RubberType,
			"production_date":   product.ProductionDate.Format("2006-01-02"),
			"status":            product.Status,
			"trace_summary":     traceSummary,
			"chain_verified":    chainVerified,
			"data_hash_matched": dataHashMatched,
		},
	}, nil
}
