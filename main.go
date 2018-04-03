// miklly project main.go
package main

import (
	"fmt"
	"strings"
	//"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/miklly/miklly/config"
	"github.com/miklly/miklly/controller"
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

func initController(w http.ResponseWriter, r *http.Request) {
	var context controller.IController
	urlPath := strings.SplitN(r.URL.Path, "/", 3)

	switch strings.ToLower(urlPath[1]) {
	case "api":
		context = new(controller.ApiController)
	case "", "web":
		context = new(controller.WebController)
	default:
		c := new(controller.Controller)
		c.Init(w, r)
		c.Error(404, "页面未找到!!")
		return
	}

	context.Init(w, r)
}

func main() {
	time.Local, _ = time.LoadLocation("PRC")
	//go taskHandler()
	config.CheckDataBase()
	//静态文件
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./wwwroot/images"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./wwwroot/css"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./wwwroot/fonts"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./wwwroot/js"))))
	//http.HandleFunc("/api/sxjhelper/record-upload", recordUpload)
	http.HandleFunc("/", initController)
	fmt.Println("启动web服务在端口：8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("启动web服务失败：", err)
	}
}
