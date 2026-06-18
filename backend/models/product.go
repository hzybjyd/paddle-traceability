package models

import "time"

type Product struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	ProductUID        string    `gorm:"uniqueIndex;size:19;not null" json:"product_uid"`
	Brand             string    `gorm:"size:50" json:"brand"`
	Model             string    `gorm:"size:100" json:"model"`
	Material          string    `gorm:"size:100" json:"material"`
	RubberType        string    `gorm:"size:100" json:"rubber_type"`
	BatchNo           string    `gorm:"size:50" json:"batch_no"`
	ProductionDate    time.Time `json:"production_date"`
	QualityReportHash string    `gorm:"size:64" json:"quality_report_hash,omitempty"`
	FactoryID         uint      `json:"factory_id"`
	Status            string    `gorm:"size:20;default:PRODUCED" json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
