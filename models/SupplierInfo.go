/*
	供货商信息数据实体模型
*/
package models

import (
	"github.com/jinzhu/gorm"
)

//供货商实体
type SupplierInfo struct {
	gorm.Model
	//供货商名称
	Name string
	//备注
	Description string
}
