package hs

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pi12138/toolbox/internal/app/hospitalization_statistics/model"
	"gorm.io/gorm"
)

func SaveToDB(tx *gorm.DB, Body *CrawlRespBody) error {
	var daily model.Daily
	daily.Date = Body.Data.Date
	daily.Cost = int(Body.Data.DailyCost)
	if err := tx.Create(&daily).Error; err != nil {
		return err
	}

	var items []model.Item
	for _, i := range Body.Data.Items {
		items = append(items, model.Item{
			Cost:      i.Cost,
			DeptName:  i.DeptName,
			ItemName:  i.ItemName,
			ItemPrice: i.ItemPrice,
			ItemQty:   i.ItemQty,
			ItemSpecs: i.ItemSpecs,
			ItemUnits: i.ItemUnits,
			TradeTime: i.TradeTime,
			VisitId:   i.VisitId,
			DailyId:   &daily.ID,
		})
	}
	return tx.Create(&items).Error
}

func SaveToJson(Body *CrawlRespBody) error {
	itemData, err := json.Marshal(Body.Data)
	if err != nil {
		return fmt.Errorf("json.Marshal error. %w", err)
	}
	filename := fmt.Sprintf("tmp/hs/data/%s.json", Body.Data.Date)
	if err := os.WriteFile(filename, itemData, 0666); err != nil {
		return fmt.Errorf(`os.WriteFile %s error. %w`, filename, err)
	}
	return nil
}
