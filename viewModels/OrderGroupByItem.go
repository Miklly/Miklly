package viewModels

import (
	"github.com/jinzhu/gorm"
	"github.com/miklly/miklly/models"
)

type OrderGroupByItem struct {
	Image          string `json:"image"`
	ThumbnailImage string `json:"thumbnailImage"`
	ImageID        uint   `json:"imageID"`
	Supplier       string `json:"supplier"`
	Count          int    `json:"count"`
}

//获取订单商品分组
func GetOrderGroupByItem(db *gorm.DB) []OrderGroupByItem {
	var list []models.OrderInfo
	db.Where("send_time is null").Find(&list)
	var result []OrderGroupByItem
	for _, value := range list {
		for _, item := range value.Items {
			var resultItem *OrderGroupByItem
			for _, i := range result {
				if item.ImageInfoID == i.ImageID && item.SupplierInfo.Name == i.Supplier {
					resultItem = &i
					goto hasImage
				}
			}
			resultItem = new(OrderGroupByItem)
			resultItem.Image = item.ImageInfo.FilePath
			resultItem.ThumbnailImage = item.ImageInfo.ThumbnailPath
			resultItem.ImageID = item.ImageInfoID
			resultItem.Supplier = item.SupplierInfo.Name
			result = append(result, *resultItem)
		hasImage:
			resultItem.Count = resultItem.Count + 1
		}
	}
	return result
}
