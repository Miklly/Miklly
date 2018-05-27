package models

import (
	"github.com/miklly/miklly/config"
)

func init() {
	db := config.OpenDataBase()
	// 全局禁用表名复数
	db.SingularTable(true)

	//var order models.OrderInfo
	//var record models.SupplierRecord
	if !db.HasTable(&ChannelInfo{}) {
		db.CreateTable(&ChannelInfo{})
		channels := []ChannelInfo{
			ChannelInfo{WXID: "als1888", Name: "VIP"},
			ChannelInfo{WXID: "ss16998", Name: "生活馆"},
			ChannelInfo{WXID: "als2888", Name: "生活馆3"},
			ChannelInfo{WXID: "abc28899", Name: "菲菲家"},
			ChannelInfo{WXID: "as-shenghuo", Name: "爱尚生活"},
			ChannelInfo{WXID: "yap0910", Name: "玩潮流"},
		}
		for _, v := range channels {
			db.Create(&v)
		}
	}
	db.AutoMigrate(&ImageInfo{},
		//&models.ChannelInfo{},
		&ErrorLog{},
		&ImageInfo{},
		&OrderInfo{},
		&OrderItem{},
		&SupplierInfo{},
		&SupplierRecord{},
		&SupplierWX{},
	)
	//db.Model(&order).Related(&order.Items)
	//db.Model(&record).Related(&record.Images, "image_info")
	db.Close()
}
