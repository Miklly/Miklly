package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/nfnt/resize"

	"github.com/jinzhu/gorm"
	. "github.com/miklly/miklly/System/Function"
	"github.com/miklly/miklly/System/Web"
	"github.com/miklly/miklly/config"
	"github.com/miklly/miklly/models"
	"github.com/miklly/miklly/viewModels"
)

func WhereAllNotSend(db *gorm.DB) *gorm.DB {
	return db.Where("send_time is null")
}

func GetOrderByID(id int) models.OrderInfo {
	db := config.OpenDataBase()
	defer db.Close()
	var order models.OrderInfo
	db.First(&order, id)
	order.LoadAtt(db)
	return order
}
func GetViewOrderByID(id int) viewModels.OrderDetail {
	result := viewModels.OrderDetail{}
	result.Init(GetOrderByID(id))
	return result
}
func GetNotSendOrderByID(id int) models.OrderInfo {
	db := config.OpenDataBase()
	var order models.OrderInfo
	WhereAllNotSend(db).First(&order, id)
	order.LoadAtt(db)
	db.Close()
	return order
}
func GetNotSendOrderByUser(name string, phone string, address string) models.OrderInfo {
	db := config.OpenDataBase()
	var order models.OrderInfo
	WhereAllNotSend(db).Where("name=? and phone=? and address=?", name, phone, address).First(&order)
	order.LoadAtt(db)
	db.Close()
	return order
}

//GetViewOrderByImageID 根据商品图片获取订单信息
func GetViewOrderByImageID(id int) []viewModels.OrderDetail {
	var list []models.OrderInfo
	var result []viewModels.OrderDetail
	db := config.OpenDataBase()
	defer db.Close()
	WhereAllNotSend(db).Find(&list)
	for _, info := range list {
		info.LoadAtt(db)
		//fmt.Println(info)
		for _, item := range info.Items {
			if item.ImageInfoID == uint(id) {
				//fmt.Println(item.ImageInfoID)
				detail := viewModels.OrderDetail{}
				detail.Init(info)
				//fmt.Println(detail)
				result = append(result, detail)
				break
			}
		}
	}
	return result
}

//GetViewOrderByItem 获取订单商品分组
func GetViewOrderByItem() []viewModels.OrderGroupByItem {
	db := config.OpenDataBase()
	defer db.Close()
	var list []models.OrderInfo
	WhereAllNotSend(db).Find(&list)
	var result []viewModels.OrderGroupByItem
	for _, value := range list {
		value.LoadAtt(db)
		for _, item := range value.Items {
			var resultItem *viewModels.OrderGroupByItem
			for i := 0; i < len(result); i++ {
				ri := result[i]
				if item.ImageInfoID == ri.ImageID && item.SupplierInfo.Name == ri.Supplier {
					resultItem = &ri
					goto hasImage
				}
			}
			resultItem = new(viewModels.OrderGroupByItem)
			resultItem.Image = item.ImageInfo.FilePath
			resultItem.ThumbnailImage = item.ImageInfo.ThumbnailPath
			resultItem.ImageID = item.ImageInfoID
			resultItem.Supplier = item.SupplierInfo.Name
			result = append(result, *resultItem)
			resultItem = &result[len(result)-1]
		hasImage:
			resultItem.Count = resultItem.Count + 1
		}
	}
	return result
}

//GetViewOrderByUser 按渠道名分组获取用户订单
func GetViewOrderByUser() []viewModels.OrderGroupByUser {
	db := config.OpenDataBase()
	defer db.Close()
	var list []models.OrderInfo
	WhereAllNotSend(db).Find(&list)
	var result []viewModels.OrderGroupByUser
	var wxGroup *viewModels.OrderGroupByUser
	for _, value := range list {
		value.LoadAtt(db)
		v := new(viewModels.OrderGroupByUserItem)
		v.FromModel(&value)
		getGroupUserByName := func(wxName string) *viewModels.OrderGroupByUser {
			for i := 0; i < len(result); i++ {
				if wxName == result[i].WX {
					return &result[i]
				}
			}
			user := viewModels.OrderGroupByUser{
				WX: wxName,
			}
			result = append(result, user)
			return &result[len(result)-1]
		}
		wxGroup = getGroupUserByName(value.ChannelInfo.Name)
		//fmt.Println(v)
		wxGroup.Users = append(wxGroup.Users, *v)
		//fmt.Println(wxGroup.Users)
	}
	//fmt.Println(result)
	return result
}

//GetViewOrderHistory 获取订单历史记录
func GetViewOrderHistory(key string, indexPage int, pageSize int) []viewModels.OrderGroupByUserItem {
	db := config.OpenDataBase()
	defer db.Close()
	var list []models.OrderInfo
	var result []viewModels.OrderGroupByUserItem
	if indexPage < 1 {
		indexPage = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	strKey := "%" + key + "%"
	db.Order("send_time desc").Where("send_time is not null and (address like ? or express_number like ? or name like ? or phone like ?) ",
		strKey, strKey, strKey, strKey).Offset((indexPage - 1) * pageSize).Limit(pageSize).Find(&list)
	//fmt.Println(list)
	for _, value := range list {
		value.LoadAtt(db)
		v := new(viewModels.OrderGroupByUserItem)
		v.FromModel(&value)
		result = append(result, *v)
	}
	return result
}
func GetChannels() []models.ChannelInfo {
	db := config.OpenDataBase()
	var list []models.ChannelInfo
	db.Find(&list)
	return list
}

//DeleteOrderByID 删除订单
func DeleteOrderByID(id int) bool {
	order := &models.OrderInfo{}
	order.ID = uint(id)
	db := config.OpenDataBase()
	db.Delete(order)
	db.Unscoped().First(order, id)
	db.Close()

	return !order.DeletedAt.IsZero()
}
func SaveOrder(order *models.OrderInfo) {
	db := config.OpenDataBase()
	SaveOrderByDB(db, order)
	db.Close()
}
func SaveOrderByDB(db *gorm.DB, order *models.OrderInfo) {
	tx := db.Begin()
	//已发货状态的订单中存在未发货的商品.则把未发货的商品提取生成一个新订单
	if order.SendTime != nil && !order.SendTime.IsZero() {
		var notSendItems []models.OrderItem
		var sendItems []models.OrderItem
		for _, item := range order.Items {
			if item.IsSend {
				sendItems = append(sendItems, item)
			} else {
				notSendItems = append(notSendItems, item)
			}
		}
		if len(notSendItems) > 0 {
			var newOrder = &models.OrderInfo{
				Address:        order.Address,
				ChannelInfoID:  order.ChannelInfoID,
				ExpressCompany: order.ExpressCompany,
				ExpressNumber:  order.ExpressNumber,
				Items:          notSendItems,
				Name:           order.Name,
				Phone:          order.Phone,
			}
			//newOrder.ID = order.ID
			//newOrder.Items = notSendItems
			order.Items = sendItems
			tx.Create(newOrder)
			newOrder.Items = notSendItems
			for i := 0; i < len(notSendItems); i++ {
				newOrder.Items[i].OrderInfoID = newOrder.ID
			}
			//tx.Update(newOrder)
			tx.Update(newOrder.Items)
			if config.HasErr(tx.Error, "创建OrderInfo记录失败!", newOrder) {
				tx.Rollback()
				return
			}
		}
	}
	tx.Save(order)
	if config.HasErr(tx.Error, "保存更新OrderItem记录失败!", order) {
		tx.Rollback()
		return
	}

	//删除不在items中的商品项
	ids := []uint{}
	for _, v := range order.Items {
		//fmt.Println(v.ID)
		ids = append(ids, v.ID)
	}
	//tx.Where("order_info_id = ? and id not in(?)", order.ID, ids).Delete(models.OrderItem{})
	tx.Delete(models.OrderItem{}, "order_info_id = ? and id not in(?)", order.ID, ids)
	if config.HasErr(tx.Error, "删除已发货Item记录失败!", ids) {
		tx.Rollback()
		return
	}
	tx.Commit()
}

//UpdateImageByString 根据图片的文本形式获取或插入图片
func UpdateImageByString(strImage string) models.ImageInfo {
	b64 := strImage[strings.Index(strImage, ",")+1:]
	data, _ := base64.StdEncoding.DecodeString(b64)
	length := len(data)
	md5 := config.MD5(data)
	db := config.OpenDataBase()
	defer db.Close()
	var result models.ImageInfo
	db.Where("length=? and m_d5=?", length, md5).First(&result)
	if result.ID < 1 {
		result.Length = uint(length)
		result.MD5 = md5
		db.Create(&result)
		fileExtend := strImage[strings.Index(strImage, "/")+1 : strings.Index(strImage, ";")]
		result.FilePath = fmt.Sprintf("/images/wx/%d.%s", result.ID, fileExtend)
		result.ThumbnailPath = fmt.Sprintf("/images/wx/%d-thumbnail.%s", result.ID, fileExtend)

		dir, _ := os.Getwd()
		dir = StrAdd(dir, "/", Web.App.Configs.WebRoot)
		f, _ := os.Create(StrAdd(dir, result.FilePath))
		f.Write(data)
		f.Sync()
		f.Close()

		dr := bytes.NewReader(data)
		var img image.Image
		switch fileExtend {
		case "png":
			img, _ = png.Decode(dr)
		case "jpeg", "jpg":
			img, _ = jpeg.Decode(dr)
		}
		bound := img.Bounds()
		imgX := float64(bound.Dx())
		imgY := float64(bound.Dy())
		var width, height uint
		if (imgX / imgY) > config.ThumbnailScale {
			width = config.ThumbnailWidth
			height = uint(imgY * (float64(config.ThumbnailWidth) / imgX))
		} else {
			height = config.ThumbnailHeight
			width = uint(imgX * (float64(config.ThumbnailHeight) / imgY))
		}

		dst := resize.Thumbnail(width, height, img, resize.Lanczos3)
		f, _ = os.Create(StrAdd(dir, result.ThumbnailPath))
		png.Encode(f, dst)
		switch fileExtend {
		case "png":
			png.Encode(f, dst)
		case "jpeg", "jpg":
			jpeg.Encode(f, dst, nil)
		}
		f.Sync()
		f.Close()

		db.Update(&result)
	}
	return result
}

//UpdateSupplierByName 根据供应商名称获取或插入
func UpdateSupplierByName(supplierName string) models.SupplierInfo {
	db := config.OpenDataBase()
	defer db.Close()
	supplier := models.SupplierInfo{}
	db.Where("name = ?", supplierName).First(&supplier)
	if supplier.ID < 1 {
		supplier.Name = supplierName
		db.Create(&supplier)
	}
	return supplier
}
