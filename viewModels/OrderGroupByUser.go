package viewModels

type OrderGroupByUser struct {
	WX    string                 `json:"wx"`
	Users []OrderGroupByUserItem `json:"users"`
}
