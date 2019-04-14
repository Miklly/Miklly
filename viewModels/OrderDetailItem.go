package viewModels

type OrderDetailItem struct {
	ImageID      uint     `json:"imageID"`
	Url          string   `json:"url"`
	Thumb        string   `json:"thumb"`
	Size         string   `json:"size"`
	Supplier     string   `json:"supplier"`
	Time         string   `json:"time"`
	IsGet        bool     `json:"isget"`
	DetailImages []string `json:"detailImages"`
}
