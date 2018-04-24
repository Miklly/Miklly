package viewModels

import (
	"github.com/miklly/miklly/models"
)

type OrderDetail struct {
	ID       uint              `json:"id"`
	WX       string            `json:"wx"`
	Name     string            `json:"name"`
	Phone    string            `json:"phone"`
	Address  string            `json:"address"`
	Company  string            `json:"company"`
	Number   string            `json:"number"`
	SendTime string            `json:"sendTime"`
	Images   []OrderDetailItem `json:"images"`
}

func (this *OrderDetail) Init(info models.OrderInfo) {
	this.ID = info.ID
	this.WX = info.ChannelInfo.Name
	this.Name = info.Name
	this.Phone = info.Phone
	this.Address = info.Address
	this.Company = info.ExpressCompany
	this.Number = info.ExpressNumber
	if !info.SendTime.IsZero() {
		this.SendTime = info.SendTime.Format("2006-01-02")
	}
	this.Images = make([]OrderDetailItem, len(info.Items))
	for index, item := range info.Items {
		this.Images[index] = OrderDetailItem{
			IsGet:        !item.GetTime.IsZero(),
			ImageID:      item.ImageInfoID,
			Size:         item.Size,
			Supplier:     item.SupplierInfo.Name,
			Time:         item.CreatedAt.Format("2006-01-02"),
			Url:          item.ImageInfo.FilePath,
			Thumb:        item.ImageInfo.ThumbnailPath,
			DetailImages: []string{item.ImageInfo.FilePath},
		}
	}
}
