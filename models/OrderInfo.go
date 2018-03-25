/*
	订单信息实体模型
*/
package models

import "time"

//订单信息实体
type OrderInfo struct {
	ID int
	//收件人
	Name string
	//联系电话
	Phone string
	//收货地址
	Address string
	//销售渠道编号
	ChannelID string
	//快递公司
	ExpressCompany string
	//快递单号
	ExpressNumber string
	//发货时间
	SendTime time.Time
}
