package mysql

import (
	"fmt"
	"log"

	"github.com/Fajar-Islami/simple_manage_products/internal/daos"
	"github.com/Fajar-Islami/simple_manage_products/internal/helper"

	"gorm.io/gorm"
)

func RunMigration(mysqlDB *gorm.DB) {
	err := mysqlDB.AutoMigrate(
		&daos.OrderHistory{},
		&daos.OrderItems{},
		&daos.User{},
	)

	if err != nil {
		log.Println(err)
	}

	var count int64
	if mysqlDB.Migrator().HasTable(&daos.OrderItems{}) {
		mysqlDB.Model(&daos.OrderItems{}).Count(&count)
		if count < 1 {
			mysqlDB.CreateInBatches(orderItemsSeed, len(orderItemsSeed))
		}
	}

	if err != nil {
		helper.Logger(currentfilepath, helper.LoggerLevelError, "", fmt.Errorf("Failed Database Migrated : %s", err.Error()))
	}

	helper.Logger(currentfilepath, helper.LoggerLevelInfo, "Database Migrated", nil)
}
