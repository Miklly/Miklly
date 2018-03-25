/*
	订单内容项数据实体模型
*/
package models

import "time"

//订单项
type OrderItem struct {
	ID int
	//订单编号
	OrderID int
	//商品图片编号
	ImageID int
	//尺寸
	Size string
	//供货商编号
	SupplierID int
	//创建时间
	CreateTime time.Time
	//拿货时间
	GetTime time.Time
	//是否发货
	IsSend bool
}
