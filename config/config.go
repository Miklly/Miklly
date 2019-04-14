package config

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const DBType string = "sqlite3"
const DBFile string = "./data.db"
const ThumbnailWidth uint = 116
const ThumbnailHeight uint = 144
const ThumbnailScale float64 = float64(29) / 36

func HasErr(err error, msg string, data ...interface{}) bool {
	if err != nil {
		Log(msg, data, err)
		return true
	}
	return false
}

var dbErrFile, _ = os.Open("Log/DBErr.log")

func OpenDataBase() *gorm.DB {
	db, err := gorm.Open(DBType, DBFile)
	HasErr(err, "打开数据库失败!", DBType, DBFile)
	//db.Set("gorm:auto_preload", true)
	db.SingularTable(true)
	db.LogMode(true)
	//db.SetLogger(gorm.Logger{level.TRACE})
	db.SetLogger(log.New(dbErrFile, "\r\n", 0))
	return db
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
