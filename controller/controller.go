package controller

import (
	"html/template"
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
)

type IController interface {
	Init(http.ResponseWriter, *http.Request)
}
type Controller struct {
	Response http.ResponseWriter
	Request  *http.Request
	db       *gorm.DB
	Action   string
	Key      string
}

func (this *Controller) Init(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.SplitN(r.URL.Path, "/", 4)
	this.Request = r
	this.Response = w
	if len(urlPath) > 2 {
		this.Action = strings.ToLower(urlPath[2])
		if len(urlPath) > 3 {
			this.Key = strings.ToLower(urlPath[3])
		}
	}
}
func (this *Controller) Error(code int, msg string) {
	this.Response.WriteHeader(code)
	t, _ := template.ParseFiles("templates/Err.gohtml")
	var data struct {
		Code int
		Msg  string
	}
	data.Code = code
	data.Msg = msg
	t.Execute(this.Response, data)
}
