package models

import "time"

type LogisticsRecord struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	ProductID     uint      `gorm:"index" json:"product_id"`
	Action        string    `gorm:"size:20" json:"action"`
	OperatorID    uint      `json:"operator_id"`
	WarehouseName string    `gorm:"size:100" json:"warehouse_name"`
	Location      string    `gorm:"size:200" json:"location"`
	Carrier       string    `gorm:"size:100" json:"carrier"`
	Remark        string    `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time `json:"created_at"`
}
