package services

import (
	"fmt"
	"main/domain"
	"os"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDb(fromScratch bool) {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}
	dsn := "root:secret@tcp(" + dbHost + ":3306)/recommender?charset=utf8mb4&parseTime=True&loc=Local"
	var e1 error
	retries := 0
	for e1 != nil || DB == nil {
		DB, e1 = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if e1 != nil {
			fmt.Println("Trying to connect to db")
			retries++
			time.Sleep(5 * time.Second)
		}
	}
	fmt.Println("Connection to db established")
	e2 := DB.AutoMigrate(&domain.Rule{}, &domain.ItemCosine{})
	HandleError(errors.Wrap(e2, "db.AutoMigrate failed"), nil, true)

	if fromScratch {
		TruncateTable("rules")
		TruncateTable("item_cosines")
	}
}

func InsertOnDb(v interface{}) {
	DB.CreateInBatches(v, 100)
}

func TruncateTable(tableName string) {
	DB.Exec("truncate table " + tableName)
}
