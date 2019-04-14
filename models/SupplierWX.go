/*
	供货商微信数据实体模型
*/
package models

import (
	"github.com/jinzhu/gorm"
)

//供货商微信
type SupplierWX struct {
	gorm.Model
	//微信号
	WXID string
	//供货商编号
	SupplierInfoID uint
	SupplierInfo   SupplierInfo
	//备注
	Description string
}

func (this *SupplierWX) LoadAtt(db *gorm.DB) {
	db.Model(this).Related(&this.SupplierInfo)
}
