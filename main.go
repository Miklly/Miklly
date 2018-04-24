// miklly project main.go
package main

import (
	"time"

	"github.com/miklly/miklly/controller"
	mvc "github.com/miklly/wemvc"
)

func main() {
	time.Local, _ = time.LoadLocation("PRC")
	//go taskHandler()
	//config.CheckDataBase()
	//静态文件
	//http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./wwwroot/images"))))
	//http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./wwwroot/css"))))
	//http.Handle("/fonts/", http.StripPrefix("/fonts/", http.FileServer(http.Dir("./wwwroot/fonts"))))
	//http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./wwwroot/js"))))
	//http.HandleFunc("/api/sxjhelper/record-upload", recordUpload)
	mvc.Route("/web/<action>", controller.WebController{}, "")
	mvc.Route("/web/<action>/<key>", controller.WebController{}, "index")
	mvc.Route("/api/<action>", controller.ApiController{}, "groupbyuser")
	mvc.Route("/api/<action>/<key>", controller.ApiController{}, "groupbyuser")
	mvc.ServeDir("wwwroot/", "/")

	mvc.Run(8080)
}
