// miklly project main.go
package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"

	. "github.com/miklly/miklly/System/Routing"
	mvc "github.com/miklly/miklly/System/Web"
	"github.com/miklly/miklly/config"
	_ "github.com/miklly/miklly/controller"
)

//注册路由,路由注册遵循以下规则，特殊路由在先，最后是标准路由
func init() {
	//Admin域的标准路由
	RouteTable.AddRote(&RouteItem{
		Name:     "admin_area",
		Url:      "admin/{controller}/{action}",
		Defaults: map[string]interface{}{"controller": "home", "action": "index", "area": "admin"},
	})
	//标准路由
	RouteTable.AddRote(&RouteItem{
		Name:        "default",
		Url:         "{controller}/{action}/{id}",
		Defaults:    map[string]interface{}{"controller": "web", "action": "index", "id": ""},
		Constraints: make(map[string]string)})
}
func main() {
	time.Local, _ = time.LoadLocation("PRC")
	config.CheckDataBase()
	//程序意外退时，记录错误日志
	defer func() {
		if e := recover(); e != nil {
			err := e.(error)
			mvc.App.Log.Add(err.Error() + "\r\n" + string(debug.Stack()))
			fmt.Println(err)
		}
	}()
	//设置最大可同时执行的进程数
	runtime.GOMAXPROCS(runtime.NumCPU()*2 - 1)
	//监听http请求
	err := mvc.App.Run()
	fmt.Println(err)
}
