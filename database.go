package FTC_App

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initializeDatabase() (*gorm.DB, error) {
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func migrateDatabase(db *gorm.DB) error {
	// Define database table schemas using GORM's migrations
	// ...
	// Run the migrations
	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	return nil
}
