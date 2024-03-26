// main.go
package main

import (
	"encoding/csv"
	"fmt"
	"kelasbeta/miniproject3/config"
	"kelasbeta/miniproject3/db"
	"kelasbeta/miniproject3/model"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	conf := config.LoadConfig()

	db, err := db.ConnectDB(conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName)
	if err != nil {
		panic("Failed to connect to database")
	}

	if err := model.AutoMigrate(db); err != nil {
		panic("Failed to migrate database")
	}


	book := model.Book{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ISBN:      "123456789",
		Penulis:   "John Doe",
		Tahun:     2022,
		Judul:     "Sample Book",
		Gambar:    "sample.jpg",
		Stok:      5,
	}
	db.Create(&book)

	var bookToUpdate model.Book
	db.First(&bookToUpdate, "judul = ?", "Sample Book")
	bookToUpdate.Penulis = "Jane Doe"
	db.Save(&bookToUpdate)

	db.Delete(&bookToUpdate)

	var books []model.Book
	db.Find(&books)
	fmt.Println("List of books:")
	for _, b := range books {
		fmt.Printf("%d | %s | %s | %d | %s\n", b.ID, b.Judul, b.Penulis, b.Tahun, b.Gambar)
	}

	importDataFromCSV(db, "sample_books.csv")
}

func importDataFromCSV(db *gorm.DB, csvFilePath string) {
	file, err := os.Open(csvFilePath)
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV records:", err)
		return
	}

	for _, record := range records {
		book := model.Book{}
		book.ISBN = record[1]
		book.Penulis = record[2]
		
		tahun, err := strconv.ParseUint(record[3], 10, 32)
		if err != nil {
			fmt.Println("Error parsing Tahun:", err)
			continue
		}
		book.Tahun = uint(tahun)

		book.Judul = record[4]
		book.Gambar = record[5]
		
		stok, err := strconv.ParseUint(record[6], 10, 32)
		if err != nil {
			fmt.Println("Error parsing Stok:", err)
			continue
		}
		book.Stok = uint(stok)

		book.CreatedAt = time.Now()
		book.UpdatedAt = time.Now()

		var existingBook model.Book
		result := db.Where("isbn = ?", book.ISBN).First(&existingBook)
		if result.Error == nil {
			existingBook.Penulis = book.Penulis
			existingBook.Tahun = book.Tahun
			existingBook.Judul = book.Judul
			existingBook.Gambar = book.Gambar
			existingBook.Stok = book.Stok
			db.Save(&existingBook)
		} else if result.Error == gorm.ErrRecordNotFound {
			db.Create(&book)
		} else {
			fmt.Println("Error querying database:", result.Error)
		}
	}
}


