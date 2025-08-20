package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"excel-seeder/config"
	"excel-seeder/database"
	"excel-seeder/excel"
	"excel-seeder/models"
)

func main() {
	// Command line flags
	var (
		configPath = flag.String("config", "config.local.yaml", "Path to config file")
		excelPath  = flag.String("excel", "file/MasterBarang.xlsx", "Path to Excel file")
		outputMode = flag.String("output", "database", "Output mode: 'database' for direct insert, 'seeder' for SQL file generation")
		seederPath = flag.String("seeder-path", "seeder/seeder.sql", "Path for generated seeder file (when output=seeder)")
	)
	flag.Parse()

	log.Printf("Starting Excel to PostgreSQL parser...")
	log.Printf("Config: %s", *configPath)
	log.Printf("Excel: %s", *excelPath)
	log.Printf("Output mode: %s", *outputMode)

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Configuration loaded successfully")

	// Parse Excel file
	log.Printf("Parsing Excel file: %s", *excelPath)
	items, err := excel.ParseExcelToMItems(*excelPath)
	if err != nil {
		log.Fatalf("Failed to parse Excel file: %v", err)
	}
	log.Printf("Successfully parsed %d items from Excel", len(items))

	if len(items) == 0 {
		log.Printf("No items found in Excel file")
		return
	}

	// Handle output based on mode
	switch *outputMode {
	case "database":
		// Direct database insertion
		log.Printf("Connecting to database...")
		db, err := database.ConnectDB(cfg)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		defer db.Close()
		log.Printf("Database connection established")

		log.Printf("Starting batch insert to database...")
		err = models.InsertMItems(db, items)
		if err != nil {
			log.Fatalf("Failed to insert items: %v", err)
		}
		log.Printf("Successfully inserted %d items to database", len(items))

	case "seeder":
		// Generate SQL seeder file
		log.Printf("Generating SQL seeder file...")

		seederDir := filepath.Dir(*seederPath)
		err := createDirIfNotExists(seederDir)
		if err != nil {
			log.Fatalf("Failed to create seeder directory: %v", err)
		}

		err = models.GenerateSeederSQL(items, *seederPath)
		if err != nil {
			log.Fatalf("Failed to generate seeder file: %v", err)
		}
		log.Printf("Successfully generated seeder file: %s", *seederPath)
		log.Printf("You can run the seeder with: psql -d your_database -f %s", *seederPath)

	default:
		log.Fatalf("Invalid output mode: %s. Use 'database' or 'seeder'", *outputMode)
	}

	log.Printf("Process completed successfully!")
}

func createDirIfNotExists(dir string) error {
	if dir == "" || dir == "." {
		return nil
	}

	_, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("invalid directory path: %v", err)
	}

	// Check if directory exists
	if _, err := filepath.Abs(dir); err == nil {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			return nil
		})
		if err != nil {
			return os.MkdirAll(dir, 0755)
		}
	}

	return nil
}
