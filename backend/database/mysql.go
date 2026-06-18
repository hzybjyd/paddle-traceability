package database

import (
	"context"
	"fmt"
	"log"

	"paddle-traceability/config"
	"paddle-traceability/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// charsetSetting flag to prevent infinite recursion when SET NAMES is running
var charsetSetting bool

func InitDB(cfg *config.DBConfig) {
	// Explicitly specify utf8mb4 charset and utf8mb4_unicode_ci collation
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True&loc=Local&interpolateParams=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("connect database failed: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("get database instance failed: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Critical fix: MySQL container's character_set_client defaults to latin1,
	// the DSN charset parameter is not actually used by the driver.
	// Use callback to set charset before each operation, avoiding recursion via flag.
	_ = DB.Callback().Create().Before("gorm:before_create").Register("set_charset_create", setCharset)
	_ = DB.Callback().Update().Before("gorm:before_update").Register("set_charset_update", setCharset)
	_ = DB.Callback().Delete().Before("gorm:before_delete").Register("set_charset_delete", setCharset)
	_ = DB.Callback().Raw().Before("gorm:raw").Register("set_charset_raw", setCharset)
	_ = DB.Callback().Query().Before("gorm:query").Register("set_charset_query", setCharset)
	_ = DB.Callback().Row().Before("gorm:row").Register("set_charset_row", setCharset)
	log.Println("registered charset SET NAMES utf8mb4 callback")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.LogisticsRecord{},
		&models.TxRecord{},
	)
	if err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	log.Println("database connected, auto migration completed")
}

// setCharset runs SET NAMES utf8mb4 before each database operation.
// Bypasses the callback chain to avoid infinite recursion.
func setCharset(db *gorm.DB) {
	if db == nil {
		return
	}
	if charsetSetting {
		return
	}
	charsetSetting = true
	defer func() { charsetSetting = false }()

	sqlDB, err := DB.DB()
	if err != nil {
		return
	}
	conn, err := sqlDB.Conn(context.Background())
	if err != nil {
		return
	}
	defer conn.Close()
	_, _ = conn.ExecContext(context.Background(), "SET NAMES utf8mb4 COLLATE utf8mb4_unicode_ci")
}
