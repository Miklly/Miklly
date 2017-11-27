// miklly project main.go
package main

import (
	"fmt"
	"log"
	"time"

	"gitee.com/johng/gkvdb/gkvdb"
)

func main() {
	time.Local, _ = time.LoadLocation("PRC")
	t := time.Now()
	fmt.Println(t.Zone())
	fmt.Println("Hello World!", t.Format("2006-01-02 15:04:05"))

	// 创建数据库，指定数据库存放目录，数据库名称
	db, err := gkvdb.New("data", "test")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// 插入数据
	key := []byte("key")
	value := []byte(t.Format("2006-01-02 15:04:05"))
	if err := db.Set(key, value); err != nil {
		fmt.Println(err)
	}

	// 查询数据
	key = []byte("name")
	fmt.Println(string(db.Get(key)))

	// 删除数据
	//	key = []byte("name")
	//	if err := db.Remove(key); err != nil {
	//		fmt.Println(err)
	//	}
}
func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
