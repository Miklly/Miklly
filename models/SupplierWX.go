/*
	供货商微信数据实体模型
*/
package models

//供货商微信
type SupplierWX struct {
	ID int
	//微信号
	WXID string
	//供货商编号
	SupplierID int
	//备注
	Description string
}
