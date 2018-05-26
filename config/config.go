package config

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

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
func HasDBErr(db *gorm.DB, msg string, data ...interface{}) bool {
	if db.Error != nil {
		Log(msg, data, db.Error)
		return true
	}
	return false
}
func OpenDataBase() *gorm.DB {
	db, err := gorm.Open(DBType, DBFile)
	CheckErr(err)
	//db.Set("gorm:auto_preload", true)
	db.SingularTable(true)
	return db
}
func CheckDataBase() {
	db := OpenDataBase()
	// 全局禁用表名复数
	db.SingularTable(true)

	//var order models.OrderInfo
	//var record models.SupplierRecord
	channels := []models.ChannelInfo{
		models.ChannelInfo{WXID: "als1888", Name: "VIP"},
		models.ChannelInfo{WXID: "ss16998", Name: "生活馆"},
		models.ChannelInfo{WXID: "als2888", Name: "生活馆3"},
		models.ChannelInfo{WXID: "abc28899", Name: "菲菲家"},
		models.ChannelInfo{WXID: "as-shenghuo", Name: "爱尚生活"},
		models.ChannelInfo{WXID: "yap0910", Name: "玩潮流"},
	}
	if !db.HasTable(&models.ChannelInfo{}) {
		db.CreateTable(&models.ChannelInfo{})
		for _, v := range channels {
			db.Create(&v)
		}
	}
	db.AutoMigrate(&models.ImageInfo{},
		//&models.ChannelInfo{},
		&models.ErrorLog{},
		&models.ImageInfo{},
		&models.OrderInfo{},
		&models.OrderItem{},
		&models.SupplierInfo{},
		&models.SupplierRecord{},
		&models.SupplierWX{},
	)
	//db.Model(&order).Related(&order.Items)
	//db.Model(&record).Related(&record.Images, "image_info")
	db.Close()
}

func StrAdd(arr ...string) string {
	if len(arr) == 0 {
		return ""
	}
	return strings.Join(arr, "")
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
func Log(msg string, data ...interface{}) {
	pc, file, line, _ := runtime.Caller(1)
	fName := runtime.FuncForPC(pc).Name()
	//2018-5-22 11:11:30 ==> 未找到指定的文件!
	//D:\tanmi\Documents\projects\go\src\github.com\miklly\miklly\config\config.go(99) -> Log(msg string, data interface{})
	//data{}
	str := fmt.Sprintf("%s ==> %s\r\n%s(%d) -> %s\r\n%#v",
		time.Now().Format("2006-01-02 15:04:05"),
		msg,
		file,
		line,
		fName,
		data)
	f, _ := os.OpenFile("Log/Error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	f.WriteString(str)
	f.Close()
	//ioutil.WriteFile("Log/Error.log", []byte(str), 0666)
}
