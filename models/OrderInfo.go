/*
	订单信息实体模型
*/
package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//订单信息实体
type OrderInfo struct {
	gorm.Model
	//收件人
	Name string `gorm:"size:20"`
	//联系电话
	Phone string `gorm:"size:32"`
	//收货地址
	Address string
	//销售渠道编号
	ChannelInfoID uint
	ChannelInfo   ChannelInfo
	//快递公司
	ExpressCompany string `gorm:"size:10"`
	//快递单号
	ExpressNumber string `gorm:"size:20"`
	//发货时间
	SendTime time.Time
	Items    []OrderItem
}
