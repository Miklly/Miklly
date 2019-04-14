/*
	供货商商品记录数据实体模型
*/
package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

//供货商商品记录实体
type SupplierRecord struct {
	gorm.Model
	//供货商微信编号
	SupplierWXID uint
	SupplierWX   SupplierWX
	//发布时间
	Time time.Time
	//文字标题
	Title string
	//图片编号列表
	Images []ImageInfo `gorm:"many2many:record_images;"`
}

func (this *SupplierRecord) LoadAtt(db *gorm.DB) {
	db.Model(this).Related(&this.SupplierWX)
	db.Model(this).Related(&this.Images)
}
