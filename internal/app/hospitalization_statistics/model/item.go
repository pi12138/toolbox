package model

type Base struct {
	ID        uint
	CreatedAt int
	UpdatedAt int
}

type Item struct {
	Base
	Cost      int    `json:"cost"`
	DeptName  string `json:"deptName" gorm:"type:varchar(101)"`
	ItemName  string `json:"itemName"`
	ItemPrice string `json:"itemPrice"`
	ItemQty   string `json:"itemQty"`
	ItemSpecs string `json:"itemSpecs"`
	ItemUnits string `json:"itemUnits"`
	TradeTime string `json:"tradeTime"`
	VisitId   string `json:"visitId"`
	DailyId   *uint
}

func (i Item) TableName() string {
	return "hs_item"
}
