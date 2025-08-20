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

// ExcelHeaderMapping mapping header Excel ke field struct (case-insensitive)
var ExcelHeaderMapping = map[string]string{
	"kode barang":    "Barcode",
	"nama barang":    "ItemName",
	"harga beli":     "PriceBase",
	"harga jual":     "DefaultPriceSale",
	"jumlah partai1": "WholesaleMinQty",
	"harga partai1":  "WholesaleUnitPrice",
	"jumlah partai2": "Wholesale2MinQty",
	"harga partai2":  "Wholesale2UnitPrice",
}

// RequiredFields daftar field yang wajib diisi
var RequiredFields = map[string]bool{
	"ItemName":  true,
	"PriceBase": true,
}

// findColumnIndex mencari index kolom berdasarkan header
func findColumnIndex(headers []string, targetHeader string) int {
	targetLower := strings.ToLower(strings.TrimSpace(targetHeader))
	for i, header := range headers {
		headerLower := strings.ToLower(strings.TrimSpace(header))
		if headerLower == targetLower {
			return i
		}
	}
	return -1
}

// ParseExcelToMItems membaca Excel dengan header mapping yang fleksibel
func ParseExcelToMItems(filename string) ([]models.MItem, error) {
	f, err := excelize.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening Excel file: %v", err)
	}
	defer f.Close()

	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("error reading Excel rows: %v", err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("Excel file is empty")
	}

	headers := rows[0]
	columnIndexes := make(map[string]int)

	for excelHeader, structField := range ExcelHeaderMapping {
		index := findColumnIndex(headers, excelHeader)
		if index != -1 {
			columnIndexes[structField] = index
			log.Printf("Mapped '%s' -> %s (column %d)", excelHeader, structField, index)
		}
	}

	var items []models.MItem
	for i, row := range rows {
		if i == 0 { // Skip header
			continue
		}

		item := models.MItem{
			IsActive:  true,
			CreatedAt: utils.TimePtr(time.Now()),
			UpdatedAt: utils.TimePtr(time.Now()),
		}

		// Set values berdasarkan column mapping
		if err := setItemValues(&item, row, columnIndexes); err != nil {
			log.Printf("Row %d: %v, skipping", i+1, err)
			continue
		}

		items = append(items, item)
	}

	return items, nil
}

// setItemValues mengatur nilai item berdasarkan column indexes
func setItemValues(item *models.MItem, row []string, columnIndexes map[string]int) error {
	// Helper function untuk get cell value
	getCellValue := func(fieldName string) string {
		if index, exists := columnIndexes[fieldName]; exists && index < len(row) {
			return strings.TrimSpace(row[index])
		}
		return ""
	}

	// Set ItemName (required)
	if itemName := getCellValue("ItemName"); itemName != "" {
		item.ItemName = itemName
	} else {
		return fmt.Errorf("ItemName is required")
	}

	// Set PriceBase (required)
	if priceStr := getCellValue("PriceBase"); priceStr != "" {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err != nil {
			return fmt.Errorf("invalid PriceBase '%s': %v", priceStr, err)
		}
		item.PriceBase = price
	} else {
		item.PriceBase = 0
		// return fmt.Errorf("PriceBase is required")
	}

	// Set Barcode (optional)
	if barcode := getCellValue("Barcode"); barcode != "" {
		item.Barcode = utils.StringPtr(barcode)
	}

	// Set DefaultPriceSale (optional)
	if defaultPriceStr := getCellValue("DefaultPriceSale"); defaultPriceStr != "" {
		if price, err := strconv.ParseFloat(defaultPriceStr, 64); err == nil {
			item.DefaultPriceSale = utils.Float64Ptr(price)
		} else {
			log.Printf("Warning: invalid DefaultPriceSale '%s', skipping", defaultPriceStr)
		}
	}

	// Set Wholesale Tier 1 fields
	if wholesaleMinQtyStr := getCellValue("WholesaleMinQty"); wholesaleMinQtyStr != "" {
		if qty, err := strconv.ParseFloat(wholesaleMinQtyStr, 64); err == nil {
			item.WholesaleMinQty = utils.Float64Ptr(qty)
		}
	}

	if wholesalePriceStr := getCellValue("WholesaleUnitPrice"); wholesalePriceStr != "" {
		if price, err := strconv.ParseFloat(wholesalePriceStr, 64); err == nil {
			item.WholesaleUnitPrice = utils.Float64Ptr(price)
		}
	}

	// Set Wholesale Tier 2 fields
	if wholesale2MinQtyStr := getCellValue("Wholesale2MinQty"); wholesale2MinQtyStr != "" {
		if qty, err := strconv.ParseFloat(wholesale2MinQtyStr, 64); err == nil {
			item.Wholesale2MinQty = utils.Float64Ptr(qty)
		}
	}

	if wholesale2PriceStr := getCellValue("Wholesale2UnitPrice"); wholesale2PriceStr != "" {
		if price, err := strconv.ParseFloat(wholesale2PriceStr, 64); err == nil {
			item.Wholesale2UnitPrice = utils.Float64Ptr(price)
		}
	}

	// Set optional fields (existing)
	if code := getCellValue("Code"); code != "" {
		item.Code = utils.StringPtr(code)
	}
	if unit := getCellValue("Unit"); unit != "" {
		item.Unit = utils.StringPtr(unit)
	}
	if mnfct := getCellValue("Mnfct"); mnfct != "" {
		item.Mnfct = utils.StringPtr(mnfct)
	}
	if spec := getCellValue("Spec"); spec != "" {
		item.Spec = utils.StringPtr(spec)
	}
	if weightStr := getCellValue("Weight"); weightStr != "" {
		if weight, err := strconv.ParseFloat(weightStr, 64); err == nil {
			item.Weight = utils.Float64Ptr(weight)
		} else {
			log.Printf("Warning: invalid Weight '%s', skipping", weightStr)
		}
	}

	return nil
}
