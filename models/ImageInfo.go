/*
	图片媒体数据实体模型
*/
package models

import "github.com/jinzhu/gorm"

//图片媒体实体
type ImageInfo struct {
	gorm.Model
	//文件大小(字节)
	Length uint
	//文件md5值(32位)
	MD5 string `gorm:"size:32"`
	//文件sha1摘要值
	SHA1 string `gorm:"size:32"`
	//文件存放路径
	FilePath string
	//缩略图存放路径
	ThumbnailPath string
}
