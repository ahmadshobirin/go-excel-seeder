package excel

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"excel-seeder/models"
	"excel-seeder/utils"

	"github.com/xuri/excelize/v2"
)

// ParseExcelToMItems membaca Excel dan mengkonversi ke slice MItem
func ParseExcelToMItems(filename string) ([]models.MItem, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening Excel file: %v", err)
	}
	defer f.Close()

	// Baca data dari sheet pertama
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("error reading Excel rows: %v", err)
	}

	var items []models.MItem
	for i, row := range rows {
		if i == 0 { // Skip header
			continue
		}

		// Pastikan row memiliki minimal kolom yang diperlukan
		if len(row) < 4 {
			log.Printf("Row %d: insufficient columns, skipping", i+1)
			continue
		}

		// Mapping kolom Excel ke struct MItem
		item := models.MItem{
			IsActive:  true, // Default value
			CreatedAt: utils.TimePtr(time.Now()),
			UpdatedAt: utils.TimePtr(time.Now()),
		}

		// Kolom 0: Code
		if len(row) > 0 && strings.TrimSpace(row[0]) != "" {
			item.Code = utils.StringPtr(strings.TrimSpace(row[0]))
		}

		// Kolom 1: ItemName (required)
		if len(row) > 1 && strings.TrimSpace(row[1]) != "" {
			item.ItemName = strings.TrimSpace(row[1])
		} else {
			log.Printf("Row %d: item_name is required, skipping", i+1)
			continue
		}

		// Kolom 2: Unit
		if len(row) > 2 && strings.TrimSpace(row[2]) != "" {
			item.Unit = utils.StringPtr(strings.TrimSpace(row[2]))
		}

		// Kolom 3: PriceBase (required)
		if len(row) > 3 && strings.TrimSpace(row[3]) != "" {
			price, err := strconv.ParseFloat(strings.TrimSpace(row[3]), 64)
			if err != nil {
				log.Printf("Row %d: invalid price_base '%s', skipping", i+1, row[3])
				continue
			}
			item.PriceBase = price
		} else {
			log.Printf("Row %d: price_base is required, skipping", i+1)
			continue
		}

		// Kolom 4: Manufacturer (opsional)
		if len(row) > 4 && strings.TrimSpace(row[4]) != "" {
			item.Mnfct = utils.StringPtr(strings.TrimSpace(row[4]))
		}

		// Kolom 5: Spec (opsional)
		if len(row) > 5 && strings.TrimSpace(row[5]) != "" {
			item.Spec = utils.StringPtr(strings.TrimSpace(row[5]))
		}

		// Kolom 6: Barcode (opsional)
		if len(row) > 6 && strings.TrimSpace(row[6]) != "" {
			item.Barcode = utils.StringPtr(strings.TrimSpace(row[6]))
		}

		items = append(items, item)
	}

	return items, nil
}
