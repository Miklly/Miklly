package controller

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/miklly/miklly/config"
	"github.com/miklly/miklly/models"
	"github.com/miklly/miklly/viewModels"
)

type WebController struct {
	Controller
}

func (this *WebController) Init(w http.ResponseWriter, r *http.Request) {
	this.Controller.Init(w, r)
	this.Response.Header().Set("content-type", "text/html; charset=utf-8")
	db, err := gorm.Open(config.DBType, config.DBFile)
	config.CheckErr(err)
	this.db = db
	defer db.Close()

	switch fmt.Sprintf("%s%s", this.Request.Method, this.Action) {
	case "GET", "GETindex":
		this.GetIndex()
	case "GETgroupbyuser":
		this.GetGroupByUser()
	case "GETadd":
		this.GetAdd()
	case "GETdetail":
		this.GetDetail()
	case "POSTedit":
		this.PostEdit()
	default:
		this.Error(404, "页面未找到!!")
	}
}
func (this *WebController) GetIndex() {
	//this.view("index", nil)
	b, err := ioutil.ReadFile("./templates/Index.html")
	config.CheckErr(err)
	this.Response.Write(b)
}
func (this *WebController) GetGroupByUser() {
	this.view("GroupByUser", viewModels.GetOrderGroupByUser(this.db))
}
func (this *WebController) GetAdd() {
	this.view("Detail", nil)
}
func (this *WebController) GetDetail() {
	var order models.OrderInfo
	id, _ := strconv.ParseUint(this.Key, 10, 32)
	this.db.First(&order, id)
	this.view("Detail", order)
}
func (this *WebController) PostEdit() {

}
func (this *WebController) view(name string, data interface{}) {
	t := template.New(fmt.Sprintf("%s.gohtml", name))

	funcChannels := func() []models.ChannelInfo {
		var list []models.ChannelInfo
		this.db.Find(&list)
		return list
	}
	t.Funcs(template.FuncMap{"channels": funcChannels})

	file := fmt.Sprintf("templates/%s.gohtml", name)
	t, _ = t.ParseFiles(file)
	t.Execute(this.Response, data)
}
