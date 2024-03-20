package hs

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/olekukonko/tablewriter"
	"github.com/pi12138/toolbox/internal/app/hospitalization_statistics/model"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type Unit struct {
	Cost      int    `json:"cost"`
	DeptName  string `json:"deptName"`
	ItemName  string `json:"itemName"`
	ItemPrice string `json:"itemPrice"`
	ItemQty   string `json:"itemQty"`
	ItemSpecs string `json:"itemSpecs"`
	ItemUnits string `json:"itemUnits"`
	TradeTime string `json:"tradeTime"`
	VisitId   string `json:"visitId"`
}

type StatData struct {
	ItemName  string
	ItemPrice int
	ItemQty   int
	ItemUnits string
	Cost      int
	List      []Unit
}

func Atoi(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return v
}

func Stat(tx *gorm.DB) []StatData {
	var items []model.Item
	tx.Find(&items)

	keyToItems := make(map[string][]model.Item)
	for i := 0; i < len(items); i++ {
		key := fmt.Sprintf("%s--[%s]", items[i].ItemName, items[i].ItemSpecs)
		keyToItems[key] = append(keyToItems[key], items[i])
	}

	ret := make([]StatData, 0, len(keyToItems))
	for k, v := range keyToItems {
		var list []Unit
		var data StatData
		data.ItemName = k
		data.ItemPrice = Atoi(v[0].ItemPrice)
		data.ItemUnits = v[0].ItemUnits
		for _, i := range v {
			list = append(list, Unit{
				Cost:      i.Cost,
				DeptName:  i.DeptName,
				ItemName:  i.ItemName,
				ItemPrice: i.ItemPrice,
				ItemQty:   i.ItemQty,
				ItemSpecs: i.ItemSpecs,
				ItemUnits: i.ItemUnits,
				TradeTime: i.TradeTime,
				VisitId:   i.VisitId,
			})
			data.ItemQty += Atoi(i.ItemQty)
			data.Cost += i.Cost
		}
		data.List = list
		ret = append(ret, data)
	}

	sort.SliceStable(ret, func(i, j int) bool {
		return ret[i].ItemQty > ret[j].ItemQty
	})
	return ret
}

func ToTable(data []StatData) {
	builder := strings.Builder{}
	table := tablewriter.NewWriter(&builder)
	table.SetHeader([]string{"名称", "单价", "数量", "单位", "总价"})
	// table.SetColMinWidth(0, 100)
	var cost int
	for _, i := range data {
		table.Append([]string{
			handlerItemName(i.ItemName),
			strconv.Itoa(i.ItemPrice),
			strconv.Itoa(i.ItemQty),
			i.ItemUnits,
			strconv.Itoa(i.Cost),
		})
		cost += i.Cost
	}
	table.SetFooter([]string{"", "", "", "共计", strconv.Itoa(cost)})
	table.Render()

	fmt.Println(builder.String())
	// f, err := os.OpenFile("tmp/hs/report.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	// n, err := f.WriteString(builder.String())
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("write %d\n", n)
}

func handlerItemName(s string) string {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}

func ToExcel(data []StatData) {
	e := excelize.NewFile()
	defer e.Close()
	excelData := [][]string{}
	excelData = append(excelData, []string{"名称", "单价", "数量", "单位", "总价"})
	var cost int
	var maxLen int
	var sheetName = "Sheet1"
	for _, i := range data {
		if v := utf8.RuneCountInString(i.ItemName); v > maxLen {
			maxLen = v
		}
		excelData = append(
			excelData,
			[]string{
				handlerItemName(i.ItemName),
				handlerPrice(i.ItemPrice),
				strconv.Itoa(i.ItemQty),
				i.ItemUnits,
				handlerPrice(i.Cost),
			})
		cost += i.Cost
	}

	fmt.Println(maxLen)
	excelData = append(excelData, []string{"", "", "", "共计", handlerPrice(cost)})
	if err := setRows(e, sheetName, 1, excelData); err != nil {
		panic(err)
	}

	if err := e.SetColWidth(sheetName, cells[0], cells[0], float64(maxLen)+20); err != nil {
		panic(err)
	}
	if err := e.SaveAs("tmp/hs/report.xlsx"); err != nil {
		panic(err)
	}
}

func handlerPrice(v int) string {
	return strconv.FormatFloat(float64(v)/100, 'f', 2, 64)
}

var cells = []string{
	"A", "B", "C", "D", "E", "F", "G",
	"H", "I", "J", "K", "L", "M", "N",
}

func setRow(f *excelize.File, sheetName string, row int, data []string) error {
	rowStr := strconv.Itoa(row)
	for i := 0; i < len(data); i++ {
		if err := f.SetCellStr(sheetName, cells[i]+rowStr, data[i]); err != nil {
			return err
		}
	}
	return nil
}

func setRows(f *excelize.File, sheetName string, startRow int, data [][]string) error {
	for i := 0; i < len(data); i++ {
		if err := setRow(f, sheetName, startRow+i, data[i]); err != nil {
			return err
		}
	}
	return nil
}
