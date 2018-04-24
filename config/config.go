package config

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/miklly/miklly/models"
)

const DBType string = "sqlite3"
const DBFile string = "./data.db"
const ThumbnailWidth uint = 116
const ThumbnailHeight uint = 144
const ThumbnailScale float64 = float64(29) / 36

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}
func OpenDataBase() *gorm.DB {
	db, err := gorm.Open(DBType, DBFile)
	CheckErr(err)
	return db
}
func CheckDataBase() {
	db := OpenDataBase()
	// 全局禁用表名复数
	db.SingularTable(true)
	db.AutoMigrate(&models.ChannelInfo{},
		&models.ImageInfo{},
		&models.OrderInfo{},
		&models.OrderItem{},
		&models.SupplierInfo{},
		&models.SupplierRecord{},
		&models.SupplierWX{})
	db.Close()
}

//获取MD5字符串
func MD5(data []byte) string {
	result := md5.Sum(data)
	return hex.EncodeToString(result[:])
}
func SHA1(data []byte) string {
	result := sha1.Sum(data)
	return hex.EncodeToString(result[:])
}
