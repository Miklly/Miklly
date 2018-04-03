package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func recordUpload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, r.Method)
	urlPath := strings.Split(r.URL.Path, "/")
	fmt.Fprintln(w, len(urlPath))
	fmt.Fprintln(w, r.URL.Path)
}

func main() {
	time.Local, _ = time.LoadLocation("PRC")
	//go taskHandler()
	//静态文件
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./images"))))
	http.HandleFunc("/", recordUpload)
	fmt.Println("启动web服务在端口：8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("启动web服务失败：", err)
	}
}
