package controller

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/miklly/miklly/viewModels"

	mvc "github.com/miklly/miklly/System/Web"
	"github.com/miklly/miklly/service"
)

type ApiController struct {
	mvc.Controller
}
type ResultStatus struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	HistoryID int    `json:"historyID"`
}

func init() {
	mvc.App.RegisterController(ApiController{})
}

//获取绝对URL地址
func (this *ApiController) getFullUrl(strUrl string) string {
	u, _ := url.Parse(strUrl)
	r := this.Request
	scheme := "http://"
	if r.TLS != nil {
		scheme = "https://"
	}
	requestURL, _ := url.Parse(strings.Join([]string{scheme, r.Host, r.RequestURI}, ""))
	return requestURL.ResolveReference(u).String()
}

//初始化控制器
func (this *ApiController) OnLoad() {
	h := this.Response.Header()
	//允许跨域访问
	h.Set("Access-Control-Allow-Origin", "*")
}

//按用户分组获取订单列表
func (this ApiController) GroupByUser() *mvc.JsonResult {
	result := service.GetViewOrderByUser()
	return this.Json(result)
}

//获取订单的商品分组
func (this ApiController) GroupByItem() *mvc.JsonResult {
	result := service.GetViewOrderByItem()
	for i := 0; i < len(result); i++ {
		result[i].Image = this.getFullUrl(result[i].Image)
		result[i].ThumbnailImage = this.getFullUrl(result[i].ThumbnailImage)
	}
	return this.Json(result)
}
func (this ApiController) OrderByID() *mvc.JsonResult {
	id, _ := strconv.Atoi(this.RouteData["id"].(string))
	order := service.GetViewOrderByID(id)
	this.converImageInDetail(&order)
	return this.Json(order)
}
func (this ApiController) converImageInDetail(item *viewModels.OrderDetail) {
	for i := 0; i < len(item.Images); i++ {
		img := item.Images[i]
		item.Images[i].Thumb = this.getFullUrl(img.Thumb)
		for n := 0; n < len(img.DetailImages); n++ {
			item.Images[i].DetailImages[n] = this.getFullUrl(img.DetailImages[n])
		}
	}
}
func (this ApiController) OrderByImage() *mvc.JsonResult {
	id, _ := strconv.Atoi(this.RouteData["id"].(string))
	result := service.GetViewOrderByImageID(id)
	for i := 0; i < len(result); i++ {
		this.converImageInDetail(&result[i])
	}
	return this.Json(result)
}
func (this ApiController) DeleteOrder() *mvc.JsonResult {
	var result ResultStatus
	id, _ := strconv.Atoi(this.RouteData["id"].(string))
	result.Success = service.DeleteOrderByID(id)
	return this.Json(result)
}

func (this ApiController) SendOrder() *mvc.JsonResult {
	form := this.Request.Form
	id, _ := strconv.Atoi(form.Get("id"))
	imgIDs := form.Get("imgIDs")
	send := form.Get("send") == "true"
	number := form.Get("number")
	t := time.Now()
	result := ResultStatus{
		Success:   true,
		HistoryID: -1,
	}
	order := service.GetOrderByID(id)
	if order.ID == 0 {
		result.Success = false
		result.Message = fmt.Sprintf("未找到相关订单(%d)!", id)
		return this.Json(result)
	}
	if send {
		order.SendTime = &t
		order.ExpressNumber = number
	}
	ids := strings.Split(imgIDs, ",")
	for key, img := range order.Items {
		item := &order.Items[key]
		strImgID := strconv.Itoa(int(img.ImageInfoID))
		for _, imgID := range ids {
			if imgID == "" {
				continue
			}
			if strImgID == imgID {
				if send {
					item.IsSend = true
				} else if img.GetTime == nil || img.GetTime.IsZero() {
					item.GetTime = &t
				}
				goto nextImg
			}
		}
		if send {
			item.IsSend = false
		} else {
			item.GetTime = nil
		}
	nextImg:
		continue
	}
	service.SaveOrder(&order)
	result.HistoryID = int(order.ID)
	return this.Json(result)
}
func (this ApiController) ChangeNumber() *mvc.JsonResult {
	id, _ := strconv.Atoi(this.RouteData["id"].(string))
	number := this.Form.String("number")
	result := ResultStatus{
		Success: true,
	}
	order := service.GetOrderByID(id)

	if order.ID == 0 {
		result.Success = false
		result.Message = fmt.Sprintf("未找到相关订单(%d)!", id)
	}
	if strings.Trim(number, " ") == "" {
		result.Success = false
		result.Message = "快递单号不能为空!"
	}
	if result.Success {
		order.ExpressNumber = number
		service.SaveOrder(&order)
	}
	return this.Json(result)
}
func (this ApiController) History() *mvc.JsonResult {
	key := this.RouteData["id"].(string)
	page := this.QueryString.Int("page")
	//page := this.Form.Int("page")
	result := service.GetViewOrderHistory(key, page, 10)
	return this.Json(result)
}
