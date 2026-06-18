package models

import "time"

type TxRecord struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProductID   uint      `gorm:"index" json:"product_id"`
	TxHash      string    `gorm:"size:128" json:"tx_hash"`
	TxType      string    `gorm:"size:20" json:"tx_type"`
	DataHash    string    `gorm:"size:64" json:"data_hash"`
	ChainStatus string    `gorm:"size:20;default:PENDING" json:"chain_status"`
	BlockHeight int64     `json:"block_height"`
	OperatorID  uint      `json:"operator_id"`
	CreatedAt   time.Time `json:"created_at"`
}
