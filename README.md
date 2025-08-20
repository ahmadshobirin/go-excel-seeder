# Excel to PostgreSQL Seeder

Sebuah tool Go yang powerful untuk mengimpor data dari file Excel (.xlsx) ke database PostgreSQL dengan dukungan batch insert yang dioptimalkan dan opsi untuk generate file SQL seeder.

## Features

- ✅ **Batch Insert Optimization**: Menggunakan perhitungan optimal batch size (32,767 / jumlah kolom)
- ✅ **Multi-Value INSERT**: Menggunakan single query dengan multiple values untuk performa maksimal
- ✅ **Dual Output Mode**: Direct database insertion atau generate SQL seeder file
- ✅ **Configuration Management**: Support YAML configuration file
- ✅ **Error Handling**: Comprehensive error handling dan logging
- ✅ **Transaction Safety**: Menggunakan database transactions untuk data integrity
- ✅ **Flexible Column Mapping**: Mudah disesuaikan dengan struktur Excel yang berbeda

## Prerequisites

- Go 1.19 atau lebih baru
- PostgreSQL database
- File Excel (.xlsx) dengan data yang akan diimpor

## Installation

1. Clone repository ini:
```bash
git clone <repository-url>
cd excel-seeder
```

2. Install dependencies:
```bash
go mod tidy
```

3. Setup database dan jalankan migration:
```bash
psql -d your_database -f db/master_item_migration.sql
```

## Configuration

Buat file `config.local.yaml` dengan konfigurasi database:

```yaml
database:
  host: localhost
  port: 5432
  user: your_username
  password: your_password
  dbname: your_database
  sslmode: disable
```

## Usage

### 1. Direct Database Insertion (Default)

Untuk insert data langsung ke database:

```bash
# Menggunakan konfigurasi dan file default
go run main.go

# Atau dengan parameter lengkap
go run main.go -output=database -config=config.local.yaml -excel=file/MasterBarang.xlsx
```

### 2. Generate SQL Seeder File

Untuk membuat file SQL seeder yang bisa dijalankan nanti:

```bash
# Generate seeder file dengan path default
go run main.go -output=seeder

# Generate seeder file dengan custom path
go run main.go -output=seeder -seeder-path=seeder/m_item_data.sql
```

### 3. Command Line Options

| Flag | Default | Description |
|------|---------|-------------|
| `-config` | `config.local.yaml` | Path ke file konfigurasi YAML |
| `-excel` | `file/MasterBarang.xlsx` | Path ke file Excel input |
| `-output` | `database` | Mode output: `database` atau `seeder` |
| `-seeder-path` | `seeder/seeder.sql` | Path untuk file seeder yang dihasilkan |

### 4. Contoh Penggunaan Lengkap

```bash
# Development - insert langsung ke database lokal
go run main.go -config=config.local.yaml -excel=data/items.xlsx

# Production - generate seeder file untuk deployment
go run main.go -output=seeder -config=config.prod.yaml -excel=data/production_items.xlsx -seeder-path=deploy/items_seeder.sql

# Testing dengan file Excel berbeda
go run main.go -excel=test_data/sample.xlsx -output=seeder
```

### 5. Menjalankan SQL Seeder File

Setelah generate file seeder, jalankan dengan:

```bash
# Jalankan seeder file ke database
psql -d your_database -f seeder/seeder.sql

# Atau dengan connection string
psql postgresql://username:password@localhost:5432/dbname -f seeder/seeder.sql
```

## Excel File Format

File Excel harus memiliki struktur kolom sebagai berikut (Sheet1):

| Column | Description | Required |
|--------|-------------|----------|
| A | Code | Optional |
| B | ItemName | **Required** |
| C | Unit | Optional |
| D | PriceBase | **Required** |
| E | Manufacturer | Optional |
| F | Spec | Optional |
| G | Barcode | Optional |

**Catatan**: Kolom ItemName dan PriceBase adalah wajib. Kolom lain bersifat opsional.

## Database Schema

Tool ini menggunakan tabel `m_item` dengan schema yang didefinisikan di `db/master_item_migration.sql`. Pastikan tabel sudah dibuat sebelum menjalankan import.

## Performance

- **Batch Size**: Otomatis dihitung berdasarkan `32,767 / 34 kolom = 963 items per batch`
- **Multi-Value INSERT**: Menggunakan single query untuk multiple rows
- **Transaction**: Setiap batch dijalankan dalam transaction terpisah
- **Memory Efficient**: Data diproses dalam batch untuk mengoptimalkan penggunaan memory

## Logging

Tool ini menyediakan logging yang comprehensive:
``
2024/01/15 10:30:00 Starting Excel to PostgreSQL parser...
2024/01/15 10:30:00 Config: config.local.yaml
2024/01/15 10:30:00 Excel: file/MasterBarang.xlsx
2024/01/15 10:30:00 Output mode: database
2024/01/15 10:30:01 Configuration loaded successfully
2024/01/15 10:30:01 Parsing Excel file: file/MasterBarang.xlsx
2024/01/15 10:30:02 Successfully parsed 5000 items from Excel
2024/01/15 10:30:02 Connecting to database...
2024/01/15 10:30:02 Database connection established
2024/01/15 10:30:02 Starting batch insert to database...
2024/01/15 10:30:02 Using batch size: 963 (calculated from 32767/34)
2024/01/15 10:30:03 Successfully inserted batch 1-963 (963 items)
2024/01/15 10:30:04 Successfully inserted batch 964-1926 (963 items)
...
2024/01/15 10:30:10 Successfully inserted 5000 items to database
2024/01/15 10:30:10 Process completed successfully!
``


## Error Handling

- **File Validation**: Memvalidasi keberadaan file Excel dan config
- **Database Connection**: Retry mechanism untuk koneksi database
- **Data Validation**: Validasi data required fields
- **Transaction Rollback**: Automatic rollback jika terjadi error dalam batch
- **Detailed Logging**: Error messages yang informatif untuk debugging

## Directory Structure
excel-seeder/
├── config/
│   └── config.go              # Configuration management
├── database/
│   └── database.go            # Database connection
├── excel/
│   └── parser.go              # Excel file parsing
├── models/
│   └── item.go                # Data models dan database operations
├── utils/
│   └── helpers.go             # Helper functions
├── db/
│   ├── master_item_migration.sql    # Database schema
│   └── master_supplier_migration.sql
├── file/
│   └── MasterBarang.xlsx      # Sample Excel file
├── seeder/                    # Generated seeder files (auto-created)
├── config.local.yaml          # Database configuration
├── go.mod
├── go.sum
├── main.go                    # Main application
└── README.md


## License
Project ini menggunakan MIT License. Lihat file `LICENSE` untuk detail lengkap.

## Troubleshooting

### Common Issues

1. **"Failed to connect to database"**
   - Pastikan PostgreSQL service berjalan
   - Periksa konfigurasi di `config.local.yaml`
   - Pastikan database dan user sudah dibuat

2. **"Failed to parse Excel file"**
   - Pastikan file Excel ada dan tidak corrupt
   - Periksa format file (harus .xlsx)
   - Pastikan Sheet1 ada dan memiliki data

3. **"Error inserting batch"**
   - Periksa apakah tabel `m_item` sudah dibuat
   - Pastikan data types sesuai dengan schema
   - Periksa constraint violations

4. **"Permission denied creating seeder directory"**
   - Pastikan aplikasi memiliki write permission
   - Buat directory secara manual jika perlu

### Debug Mode

Untuk debugging, tambahkan logging level yang lebih detail di kode atau gunakan PostgreSQL logs untuk melihat query yang dieksekusi.
