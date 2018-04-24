package viewModels

type OrderGroupByItem struct {
	Image          string `json:"image"`
	ThumbnailImage string `json:"thumbnailImage"`
	ImageID        uint   `json:"imageID"`
	Supplier       string `json:"supplier"`
	Count          int    `json:"count"`
}
