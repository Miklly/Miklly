package viewModels

import (
	"github.com/miklly/miklly/models"
)

type OrderGroupByUserItem struct {
	ID        uint   `json:"id"`
	ChannelID string `json:"channelID"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	ItemCount int    `json:"itemCount"`
	SendTime  string `json:"sendTime"`
	Number    string `json:"number"`
}

func (this *OrderGroupByUserItem) FromModel(oi *models.OrderInfo) {
	this.ID = oi.ID
	this.Address = oi.Address
	this.ChannelID = oi.ChannelInfo.WXID
	this.ItemCount = len(oi.Items)
	this.Name = oi.Name
	this.Number = oi.ExpressNumber
	this.Phone = oi.Phone
	if oi.SendTime != nil && !oi.SendTime.IsZero() {
		this.SendTime = oi.SendTime.Format("2006-01-02")
	}
}
