package config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/miklly/miklly/models"
)

const DBType string = "sqlite3"
const DBFile string = "./data.db"

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

func CheckDataBase() {
	db, err := gorm.Open(DBType, DBFile)
	CheckErr(err)
	defer db.Close()
	// 全局禁用表名复数
	db.SingularTable(true)
	db.AutoMigrate(&models.ChannelInfo{},
		&models.ImageInfo{},
		&models.OrderInfo{},
		&models.OrderItem{},
		&models.SupplierInfo{},
		&models.SupplierRecord{},
		&models.SupplierWX{})

}
