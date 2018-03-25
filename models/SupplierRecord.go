/*
	供货商商品记录数据实体模型
*/
package models

import "time"

//供货商商品记录实体
type SupplierRecord struct {
	ID int
	//供货商微信编号
	SupplierWXID int
	//发布时间
	Time time.Time
	//文字标题
	Title string
	//图片编号列表
	PicIDs string
	//创建时间
	CreateTime time.Time
}
