package controller

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/miklly/miklly/models"
	"github.com/miklly/miklly/service"
	mvc "github.com/miklly/wemvc"
)

type WebController struct {
	mvc.Controller
}

func (this *WebController) OnInit(ctx *mvc.Context) {
	this.Controller.OnInit(ctx)
	this.Response().Header().Set("content-type", "text/html; charset=utf-8")
}
func (this WebController) Index() mvc.Result {
	return this.File("./views/Index.html", "text/html")
}
func (this WebController) GroupByUser() mvc.Result {
	this.ViewData["list"] = service.GetViewOrderByUser()
	return this.View()
}
func (this WebController) GetAdd() mvc.Result {

	return this.ViewAction("Detail")
}
func (this WebController) GetDetail() mvc.Result {
	id, _ := strconv.Atoi(this.RouteData()["key"])
	this.ViewData["info"] = service.GetOrderByID(id)
	return this.View()
}
func (this WebController) PostEdit() mvc.Result {
	var entity models.OrderInfo
	order := &models.OrderInfo{}
	mvc.ModelParse(order, this.Request().Form)
	var trimStr = func(r rune) bool {
		ts := []rune(" ,.:，。：`_+=!~|")
		for _, v := range ts {
			if v == r {
				return true
			}
		}
		return false
	}
	strings.TrimFunc(order.Address, trimStr)
	strings.TrimFunc(order.Name, trimStr)
	strings.TrimFunc(order.Phone, trimStr)
	if order.ID > 0 {
		entity = service.GetNotSendOrderByID(int(order.ID))
		if entity.ID < 1 {
			return this.Error("未找到指定的订单!")
		}
	} else {
		entity = service.GetNotSendOrderByUser(order.Name, order.Phone, order.Address)
	}

	entity.Address = order.Address
	entity.ChannelInfoID = order.ChannelInfoID
	entity.ExpressCompany = order.ExpressCompany
	entity.ExpressNumber = order.ExpressNumber
	entity.Name = order.Name
	entity.Phone = order.Phone
	entity.SendTime = order.SendTime

	//移除商品
	newItem := []models.OrderItem{}
	delImageIDs := this.Request().Form.Get("imageDelete")
	for _, id := range strings.Split(delImageIDs, ",") {
		uid, err := strconv.Atoi(id)
		if err == nil && uid > 0 {
			for _, item := range entity.Items {
				if item.ImageInfoID != uint(uid) {
					newItem = append(newItem, item)
				}
			}
		}
	}
	if len(newItem) > 0 {
		entity.Items = newItem
	}

	itemCount, err := strconv.Atoi(this.Request().Form.Get("itemCount"))
	if err == nil {
		var item models.OrderItem
		for i := 0; i < itemCount; i++ {
			strImage := this.Request().Form.Get(fmt.Sprintf("hidFile-%d", i))
			if strImage == "" {
				continue
			}
			imgID, err := strconv.Atoi(strImage)
			if err == nil {
				for _, val := range entity.Items {
					if val.ImageInfoID == uint(imgID) {
						item = val
						break
					}
				}
				if item.ID < 1 {
					continue
				}
			} else {
				img := service.UpdateImageByString(strImage)
				item = models.OrderItem{
					ImageInfoID: img.ID,
					ImageInfo:   img,
				}
				entity.Items = append(entity.Items, item)
			}
			item.Size = this.Request().Form.Get(fmt.Sprintf("size-%d", i))
			supplierName := this.Request().Form.Get(fmt.Sprintf("supplier-%d", i))
			item.SupplierInfo = service.UpdateSupplierByName(supplierName)
			item.SupplierInfoID = item.SupplierInfo.ID
		}
	}
	service.SaveOrder(&entity)
	return this.Redirect("/web/index")
}
func (this WebController) Error(msg string) mvc.Result {
	this.ViewData["msg"] = msg
	return this.View()
}
