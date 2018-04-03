package controller

import (
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/miklly/miklly/config"
	"github.com/miklly/miklly/viewModels"
)

type ApiController struct {
	Controller
}

func (this *ApiController) Init(w http.ResponseWriter, r *http.Request) {
	this.Controller.Init(w, r)
	h := this.Response.Header()
	//允许跨域访问
	h.Add("Access-Control-Allow-Origin", "*")
	h.Set("content-type", "text/json; charset=utf-8")

	db, err := gorm.Open(config.DBType, config.DBFile)
	config.CheckErr(err)
	this.db = db
	defer db.Close()

	switch this.Action {
	case "groupbyuser":
		this.groupByUser()
	default:
		this.Response.WriteHeader(404)
		this.Response.Write([]byte(`{"success":false,"content":"404 Not Found"}`))
	}
}

//按用户分组获取订单列表
func (this *ApiController) groupByUser() {
	result := viewModels.GetOrderGroupByUser(this.db)
	b, err := json.Marshal(result)
	config.CheckErr(err)
	this.Response.Write(b)
}

//获取订单的商品分组
func (this *ApiController) groupByItem() {
	result := viewModels.GetOrderGroupByItem(this.db)
	b, err := json.Marshal(result)
	config.CheckErr(err)
	this.Response.Write(b)
}
