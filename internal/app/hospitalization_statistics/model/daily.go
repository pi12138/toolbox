package model

type Daily struct {
	Base

	Date string `gorm:"unique"`
	Cost int
}

func (d Daily) TableName() string {
	return "hs_daily"
}
