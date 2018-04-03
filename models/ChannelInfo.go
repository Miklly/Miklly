/*
	销售渠道实体数据模型
*/
package models

import "github.com/jinzhu/gorm"

//销售渠道
type ChannelInfo struct {
	gorm.Model
	//微信号
	WXID string `gorm:"not null;unique"`
	//销售渠道名称
	Name string `gorm:"not null;unique"`
}
