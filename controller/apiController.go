package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/miklly/miklly/service"
	mvc "github.com/miklly/wemvc"
)

type ApiController struct {
	mvc.Controller
}
type resultStatus struct {
	success   bool
	message   string
	historyID int
}

//初始化控制器
func (this *ApiController) OnInit(ctx *mvc.Context) {
	this.Controller.OnInit(ctx)

	h := this.Response().Header()
	//允许跨域访问
	h.Add("Access-Control-Allow-Origin", "*")
	h.Set("content-type", "text/json; charset=utf-8")
}

//按用户分组获取订单列表
func (this ApiController) GroupByUser() mvc.Result {
	result := service.GetViewOrderByUser()
	return this.JSON(result)
}

//获取订单的商品分组
func (this ApiController) GroupByItem() mvc.Result {
	result := service.GetViewOrderByItem()
	return this.JSON(result)
}
func (this ApiController) OrderByID() mvc.Result {
	id, _ := strconv.Atoi(this.RouteData()["key"])
	order := service.GetOrderByID(id)
	return this.JSON(order)
}
func (this ApiController) OrderByImage() mvc.Result {
	id, _ := strconv.Atoi(this.RouteData()["key"])
	result := service.GetViewOrderByImageID(id)
	return this.JSON(result)
}
func (this ApiController) DeleteOrder() mvc.Result {
	var result resultStatus
	id, _ := strconv.Atoi(this.RouteData()["key"])
	result.success = service.DeleteOrderByID(id)
	return this.JSON(result)
}

func (this ApiController) PostSendOrder() mvc.Result {
	form := this.Request().Form
	id, _ := strconv.Atoi(this.RouteData()["key"])
	imgIDs := form.Get("imgIDs")
	send, _ := strconv.ParseBool(form.Get("send"))
	number := form.Get("number")
	result := resultStatus{
		success:   true,
		historyID: -1,
	}
	order := service.GetOrderByID(id)
	if order.ID == 0 {
		result.success = false
		result.message = fmt.Sprintf("未找到相关订单(%d)!", id)
		return this.JSON(result)
	}
	if send {
		order.SendTime = time.Now()
		order.ExpressNumber = number
	}
	ids := strings.Split(imgIDs, ",")
	for _, img := range order.Items {
		strImgID := strconv.Itoa(int(img.ImageInfoID))
		for _, imgID := range ids {
			if imgID == "" {
				continue
			}
			if strImgID == imgID {
				if send {
					img.IsSend = true
				} else if img.GetTime.IsZero() {
					img.GetTime = time.Now()
				}
				goto nextImg
			}
		}
		if send {
			img.IsSend = false
		} else {
			img.GetTime = time.Time{}
		}
	nextImg:
		continue
	}
	service.SaveOrder(&order)
	result.historyID = int(order.ID)
	return this.JSON(result)
}
func (this ApiController) ChangeNumber() mvc.Result {
	id, _ := strconv.Atoi(this.RouteData()["key"])
	number := this.Request().Form.Get("number")
	result := resultStatus{
		success: true,
	}
	order := service.GetOrderByID(id)

	if order.ID == 0 {
		result.success = false
		result.message = fmt.Sprintf("未找到相关订单(%d)!", id)
	}
	if strings.Trim(number, " ") == "" {
		result.success = false
		result.message = "快递单号不能为空!"
	}
	if result.success {
		order.ExpressNumber = number
		service.SaveOrder(&order)
	}
	return this.JSON(result)
}
func (this ApiController) History() mvc.Result {
	key := this.RouteData()["key"]
	page, _ := strconv.Atoi(this.Request().Form.Get("page"))
	result := service.GetViewOrderHistory(key, page, 10)
	return this.JSON(result)
}
