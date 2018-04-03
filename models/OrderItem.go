/*
	订单内容项数据实体模型
*/
package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//订单项
type OrderItem struct {
	gorm.Model
	//订单编号
	OrderInfoID uint
	//商品图片编号
	ImageInfoID uint
	ImageInfo   ImageInfo
	//尺寸
	Size string
	//供货商编号
	SupplierInfoID uint
	SupplierInfo   SupplierInfo
	//拿货时间
	GetTime time.Time
	//是否发货
	IsSend bool
}
