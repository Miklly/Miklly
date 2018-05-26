package models

import (
	"github.com/jinzhu/gorm"
)

type ErrorLog struct {
	gorm.Model
	File    string
	Line    string
	Message string
	Data    string
}
