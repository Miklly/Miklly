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
	Name string `gorm:"size:20;index:findOrder"`
	//联系电话
	Phone string `gorm:"size:32;index:findOrder"`
	//收货地址
	Address string `gorm:"index:findOrder"`
	//销售渠道编号
	ChannelInfoID uint
	ChannelInfo   ChannelInfo
	//快递公司
	ExpressCompany string `gorm:"size:10"`
	//快递单号
	ExpressNumber string `gorm:"size:20;index"`
	//发货时间
	SendTime *time.Time
	Items    []OrderItem
}

func (this *OrderInfo) LoadAtt(db *gorm.DB) {
	db.Model(this).Related(&this.ChannelInfo)
	db.Model(this).Related(&this.Items)
	for i := 0; i < len(this.Items); i++ {
		this.Items[i].LoadAtt(db)
	}
}
