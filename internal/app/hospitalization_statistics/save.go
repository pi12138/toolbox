package hs

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/pi12138/toolbox/internal/app/hospitalization_statistics/model"
	"gorm.io/gorm"
)

func SaveToDB(tx *gorm.DB, Body *CrawlRespBody) error {
	var daily model.Daily
	if err := tx.Where(&model.Daily{
		Date: Body.Data.Date,
	}).First(&daily).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("query Daily error. %w", err)
	}
	daily.Date = Body.Data.Date
	daily.Cost = int(Body.Data.DailyCost)
	if daily.ID != 0 {
		if err := tx.Save(&daily).Error; err != nil {
			return err
		}
		var items []model.Item
		if err := tx.Where(&model.Item{
			DailyId: &daily.ID,
		}).Find(&items).Error; err != nil {
			return fmt.Errorf("query Item error. %w", err)
		}
		if len(items) > 0 {
			if err := tx.Delete(&items).Error; err != nil {
				return err
			}
		}
	} else {
		if err := tx.Create(&daily).Error; err != nil {
			return err
		}
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
