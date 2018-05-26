package controller

import (
	"fmt"
	"strconv"
	"strings"

	mvc "github.com/miklly/miklly/System/Web"
	"github.com/miklly/miklly/models"
	"github.com/miklly/miklly/service"
)

type WebController struct {
	mvc.Controller
}

func init() {
	mvc.App.RegisterController(WebController{})
}

func (this *WebController) GroupByUser() *mvc.ViewResult {
	this.ViewData["list"] = service.GetViewOrderByUser()
	return this.View()
}
func (this *WebController) Add() *mvc.ViewResult {

	return this.Detail()
}
func (this *WebController) Detail() *mvc.ViewResult {
	id, _ := strconv.Atoi(fmt.Sprintf("%v", this.RouteData["id"]))
	this.ViewData["info"] = service.GetOrderByID(id)
	this.ViewData["channels"] = service.GetChannels()
	return this.View("Detail")
}
func (this *WebController) Edit(order models.OrderInfo) *mvc.ViewResult {
	var entity models.OrderInfo
	var trimStr = func(r rune) bool {
		ts := []rune(" ,.:，。：`_+=!~|")
		for _, v := range ts {
			if v == r {
				return true
			}
		}
		return false
	}
	order.Address = strings.TrimFunc(order.Address, trimStr)
	order.Name = strings.TrimFunc(order.Name, trimStr)
	order.Phone = strings.TrimFunc(order.Phone, trimStr)
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
	delImageIDs := strings.Split(this.Form.String("imageDelete"), ",")
	if len(delImageIDs) > 0 {
		newItem := []models.OrderItem{}
		for _, item := range entity.Items {
			for _, id := range delImageIDs {
				uid, err := strconv.Atoi(id)
				if err == nil && item.ID == uint(uid) {
					goto nextItem
				}
			}
			newItem = append(newItem, item)
		nextItem:
		}
		entity.Items = newItem
	}

	itemCount, err := strconv.Atoi(this.Form.String("itemCount"))
	//如果转换成功则继续
	if err == nil {
		var item *models.OrderItem
		for i := 0; i <= itemCount; i++ {
			item = nil
			strImage := this.Form.String(fmt.Sprintf("hidFile-%d", i))
			if strImage == "" {
				continue
			}
			itemID, err := strconv.Atoi(strImage)
			if err == nil {
				for index, val := range entity.Items {
					if val.ID == uint(itemID) {
						item = &entity.Items[index]
						break
					}
				}
				if item == nil {
					continue
				}
			} else {
				img := service.UpdateImageByString(strImage)
				item = &models.OrderItem{
					ImageInfoID: img.ID,
					ImageInfo:   img,
				}
				entity.Items = append(entity.Items, *item)
				item = &entity.Items[len(entity.Items)-1]
			}
			item.Size = this.Form.String(fmt.Sprintf("size-%d", i))
			supplierName := this.Form.String(fmt.Sprintf("supplier-%d", i))
			item.SupplierInfo = service.UpdateSupplierByName(supplierName)
			item.SupplierInfoID = item.SupplierInfo.ID
		}
	}
	service.SaveOrder(&entity)
	this.Redirect("/")
	return nil
}
func (this *WebController) Error(msg string) *mvc.ViewResult {
	this.ViewData["msg"] = msg
	return this.View()
}
