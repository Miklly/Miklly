// miklly project main.go
package main

import (
	"fmt"
	//"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"time"
	//"./config"
	"./controller"
)

func recordUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		t := r.PostFormValue("time")
		txt := r.PostFormValue("text")
		count, _ := strconv.Atoi(r.PostFormValue("fileCount"))

		for i := 0; i < count; i++ {
			file, handler, err := r.FormFile("pic" + strconv.Itoa(i))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()
			//保存的图片名称后续改为数据表ＩＤ
			f, err := os.OpenFile("./images/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer f.Close()
			io.Copy(f, file)
		}
		//保存至数据库
		fmt.Fprintf(w, "已成功保存：%s %s\r\n图片总数：%d", t, txt, count)
	} else {
		fmt.Fprintln(w, "请使用ＰＯＳＴ方法！！")
	}

}

func taskHandler() {
	lastTime := int64(0)
	for {
		now := time.Now()
		if now.Unix() >= lastTime {
			err := exec.Command("am", "broadcast", "-a", "com.miklly.mywidgetapp.UPDATE").Run()
			if err != nil {
				log.Fatal(err)
			}
			next := now.Add(time.Hour * 24)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
			lastTime = next.Unix()
		}
		time.Sleep(time.Second * 10)
	}
}

func main() {
	time.Local, _ = time.LoadLocation("PRC")
	//go taskHandler()
	//静态文件
	//http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	http.HandleFunc("/api/sxjhelper/record-upload", recordUpload)
	fmt.Println("启动web服务在端口：80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal("启动web服务失败：", err)
	}
}
