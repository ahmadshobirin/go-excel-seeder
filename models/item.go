package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type MItem struct {
	ID                  int64      `db:"id"`
	MBuID               *int64     `db:"m_bu_id"`
	Code                *string    `db:"code"`
	MItemTypeID         *int64     `db:"m_item_type_id"`
	MCat1ID             *int64     `db:"m_cat1_id"`
	MCat2ID             *int64     `db:"m_cat2_id"`
	MCat3ID             *int64     `db:"m_cat3_id"`
	MCat4ID             *int64     `db:"m_cat4_id"`
	ItemName            string     `db:"item_name"`
	ItemNameLong        *string    `db:"item_name_long"`
	UnitID              *int64     `db:"unit_id"`
	Unit                *string    `db:"unit"`
	Mnfct               *string    `db:"mnfct"`
	PriceBase           float64    `db:"price_base"`
	ItemPhoto           *string    `db:"item_photo"`
	Spec                *string    `db:"spec"`
	Weight              *float64   `db:"weight"`
	WeightUnitID        *int64     `db:"weight_unit_id"`
	DimL                *float64   `db:"dim_l"`
	DimLUnitID          *int64     `db:"dim_l_unit_id"`
	DimP                *float64   `db:"dim_p"`
	DimPUnitID          *int64     `db:"dim_p_unit_id"`
	DimT                *float64   `db:"dim_t"`
	DimTUnitID          *int64     `db:"dim_t_unit_id"`
	IsActive            bool       `db:"is_active"`
	CreatorID           *int32     `db:"creator_id"`
	EditorID            *int32     `db:"editor_id"`
	CreatedAt           *time.Time `db:"created_at"`
	UpdatedAt           *time.Time `db:"updated_at"`
	IsTimbangan         *bool      `db:"is_timbangan"`
	Round               *float64   `db:"round"`
	FlagPPN             *bool      `db:"flag_ppn"`
	MSuppID             *int64     `db:"m_supp_id"`
	DefaultPriceSale    *float64   `db:"default_price_sale"`
	Barcode             *string    `db:"barcode"`
	WholesaleMinQty     *float64   `db:"wholesale_min_qty"`
	WholesaleUnitPrice  *float64   `db:"wholesale_unit_price"`
	Wholesale2MinQty    *float64   `db:"wholesale_2_min_qty"`
	Wholesale2UnitPrice *float64   `db:"wholesale_2_unit_price"`
}

const (
	PostgreSQLParamLimit = 32767 // PostgreSQL parameter limit adalah 65535, tapi kita gunakan 32767 untuk safety
	MItemColumnCount     = 34    // Jumlah kolom dalam tabel m_item (tanpa id yang auto-increment)
)

func calculateBatchSize() int {
	return PostgreSQLParamLimit / MItemColumnCount
}

func InsertMItems(db *sql.DB, items []MItem) error {
	if len(items) == 0 {
		return nil
	}

	batchSize := calculateBatchSize()
	log.Printf("Using batch size: %d (calculated from %d/%d)", batchSize, PostgreSQLParamLimit, MItemColumnCount)

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		if err := insertBatch(db, batch); err != nil {
			return fmt.Errorf("error inserting batch %d-%d: %v", i+1, end, err)
		}
		log.Printf("Successfully inserted batch %d-%d (%d items)", i+1, end, len(batch))
	}

	return nil
}

// insertBatch melakukan insert untuk satu batch menggunakan multi-value INSERT
func insertBatch(db *sql.DB, items []MItem) error {
	if len(items) == 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback()

	// Build multi-value INSERT query
	query := `INSERT INTO m_item (
		m_bu_id, code, m_item_type_id, m_cat1_id, m_cat2_id, m_cat3_id, m_cat4_id,
		item_name, item_name_long, unit_id, unit, mnfct, price_base, item_photo,
		spec, weight, weight_unit_id, dim_l, dim_l_unit_id, dim_p, dim_p_unit_id,
		dim_t, dim_t_unit_id, is_active, creator_id, editor_id, created_at,
		updated_at, is_timbangan, round, flag_ppn, m_supp_id, default_price_sale, barcode
	) VALUES `

	valuesPlaceholders := make([]string, len(items))
	args := make([]interface{}, 0, len(items)*MItemColumnCount)

	for i, item := range items {
		paramStart := i * MItemColumnCount
		placeholders := make([]string, MItemColumnCount)
		for j := 0; j < MItemColumnCount; j++ {
			placeholders[j] = fmt.Sprintf("$%d", paramStart+j+1)
		}
		valuesPlaceholders[i] = "(" + strings.Join(placeholders, ", ") + ")"

		// Add arguments in the same order as the columns
		args = append(args,
			item.MBuID, item.Code, item.MItemTypeID, item.MCat1ID, item.MCat2ID,
			item.MCat3ID, item.MCat4ID, item.ItemName, item.ItemNameLong,
			item.UnitID, item.Unit, item.Mnfct, item.PriceBase, item.ItemPhoto,
			item.Spec, item.Weight, item.WeightUnitID, item.DimL, item.DimLUnitID,
			item.DimP, item.DimPUnitID, item.DimT, item.DimTUnitID, item.IsActive,
			item.CreatorID, item.EditorID, item.CreatedAt, item.UpdatedAt,
			item.IsTimbangan, item.Round, item.FlagPPN, item.MSuppID,
			item.DefaultPriceSale, item.Barcode,
		)
	}

	query += strings.Join(valuesPlaceholders, ", ")

	_, err = tx.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing batch insert: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error committing transaction: %v", err)
	}

	return nil
}

// GenerateSeederSQL membuat file SQL seeder dari data items
func GenerateSeederSQL(items []MItem, outputPath string) error {
	if len(items) == 0 {
		return fmt.Errorf("no items to generate seeder")
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating seeder file: %v", err)
	}
	defer file.Close()

	// Write header
	_, err = file.WriteString("-- Generated seeder file for m_item table\n")
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("-- Generated at: %s\n", time.Now().Format("2006-01-02 15:04:05")))
	if err != nil {
		return err
	}
	_, err = file.WriteString(fmt.Sprintf("-- Total items: %d\n\n", len(items)))
	if err != nil {
		return err
	}

	batchSize := calculateBatchSize()
	log.Printf("Generating SQL seeder with batch size: %d", batchSize)

	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]
		if err := writeBatchSQL(file, batch, i+1); err != nil {
			return fmt.Errorf("error writing batch %d-%d: %v", i+1, end, err)
		}
	}

	log.Printf("Successfully generated seeder file: %s", outputPath)
	return nil
}

// writeBatchSQL menulis satu batch INSERT statement ke file
func writeBatchSQL(file *os.File, items []MItem, batchNum int) error {
	if len(items) == 0 {
		return nil
	}

	// Write batch comment
	_, err := file.WriteString(fmt.Sprintf("-- Batch %d (%d items)\n", batchNum, len(items)))
	if err != nil {
		return err
	}

	// Write INSERT statement
	_, err = file.WriteString(`INSERT INTO m_item (
	m_bu_id, code, m_item_type_id, m_cat1_id, m_cat2_id, m_cat3_id, m_cat4_id,
	item_name, item_name_long, unit_id, unit, mnfct, price_base, item_photo,
	spec, weight, weight_unit_id, dim_l, dim_l_unit_id, dim_p, dim_p_unit_id,
	dim_t, dim_t_unit_id, is_active, creator_id, editor_id, created_at,
	updated_at, is_timbangan, round, flag_ppn, m_supp_id, default_price_sale, barcode
) VALUES`)
	if err != nil {
		return err
	}

	for i, item := range items {
		if i > 0 {
			_, err = file.WriteString(",")
			if err != nil {
				return err
			}
		}

		_, err = file.WriteString("\n\t(")
		if err != nil {
			return err
		}

		// Format values
		values := []string{
			formatSQLValue(item.MBuID),
			formatSQLValue(item.Code),
			formatSQLValue(item.MItemTypeID),
			formatSQLValue(item.MCat1ID),
			formatSQLValue(item.MCat2ID),
			formatSQLValue(item.MCat3ID),
			formatSQLValue(item.MCat4ID),
			formatSQLValue(item.ItemName),
			formatSQLValue(item.ItemNameLong),
			formatSQLValue(item.UnitID),
			formatSQLValue(item.Unit),
			formatSQLValue(item.Mnfct),
			formatSQLValue(item.PriceBase),
			formatSQLValue(item.ItemPhoto),
			formatSQLValue(item.Spec),
			formatSQLValue(item.Weight),
			formatSQLValue(item.WeightUnitID),
			formatSQLValue(item.DimL),
			formatSQLValue(item.DimLUnitID),
			formatSQLValue(item.DimP),
			formatSQLValue(item.DimPUnitID),
			formatSQLValue(item.DimT),
			formatSQLValue(item.DimTUnitID),
			formatSQLValue(item.IsActive),
			formatSQLValue(item.CreatorID),
			formatSQLValue(item.EditorID),
			formatSQLValue(item.CreatedAt),
			formatSQLValue(item.UpdatedAt),
			formatSQLValue(item.IsTimbangan),
			formatSQLValue(item.Round),
			formatSQLValue(item.FlagPPN),
			formatSQLValue(item.MSuppID),
			formatSQLValue(item.DefaultPriceSale),
			formatSQLValue(item.Barcode),
		}

		_, err = file.WriteString(strings.Join(values, ", "))
		if err != nil {
			return err
		}

		_, err = file.WriteString(")")
		if err != nil {
			return err
		}
	}

	_, err = file.WriteString(";\n\n")
	return err
}

// formatSQLValue memformat nilai untuk SQL statement
func formatSQLValue(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return "NULL"
	case *string:
		if v == nil {
			return "NULL"
		}
		return fmt.Sprintf("'%s'", strings.ReplaceAll(*v, "'", "''"))
	case string:
		return fmt.Sprintf("'%s'", strings.ReplaceAll(v, "'", "''"))
	case *int64:
		if v == nil {
			return "NULL"
		}
		return fmt.Sprintf("%d", *v)
	case int64:
		return fmt.Sprintf("%d", v)
	case *int32:
		if v == nil {
			return "NULL"
		}
		return fmt.Sprintf("%d", *v)
	case int32:
		return fmt.Sprintf("%d", v)
	case *float64:
		if v == nil {
			return "NULL"
		}
		return fmt.Sprintf("%f", *v)
	case float64:
		return fmt.Sprintf("%f", v)
	case *bool:
		if v == nil {
			return "NULL"
		}
		if *v {
			return "true"
		}
		return "false"
	case bool:
		if v {
			return "true"
		}
		return "false"
	case *time.Time:
		if v == nil {
			return "NULL"
		}
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	default:
		return "NULL"
	}
}
