// main_test.go
package main_test

import (
	"fmt"
	"kelasbeta/miniproject3/config"
	"kelasbeta/miniproject3/db"
	"kelasbeta/miniproject3/model"
	"os"
	"testing"
	"time"

	"gorm.io/gorm"
)

var (
	testDB     *gorm.DB
	testBookID uint
)

func TestMain(m *testing.M) {
	// Load .env file
	_ = os.Setenv("DB_HOST", "localhost")
	_ = os.Setenv("DB_PORT", "3306")
	_ = os.Setenv("DB_USER", "root")
	_ = os.Setenv("DB_PASS", "")
	_ = os.Setenv("DB_NAME", "book")

	conf := config.LoadConfig()

	var err error
	testDB, err = db.ConnectDB(conf.DBHost, conf.DBPort, conf.DBUser, conf.DBPassword, conf.DBName)
	if err != nil {
		fmt.Println("Failed to connect to test database")
		os.Exit(1)
	}

	if err := model.AutoMigrate(testDB); err != nil {
		fmt.Println("Failed to migrate test database")
		os.Exit(1)
	}

	exitVal := m.Run()

	testDB.Delete(&model.Book{}, "id = ?", testBookID)

	os.Exit(exitVal)
}

func TestCRUD(t *testing.T) {
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
	testDB.Create(&book)
	testBookID = book.ID

	var bookToUpdate model.Book
	testDB.First(&bookToUpdate, "judul = ?", "Sample Books")
	bookToUpdate.Penulis = "Jane Does"
	testDB.Save(&bookToUpdate)

	var books []model.Book
	testDB.Find(&books)
	if len(books) != 1 {
		t.Errorf("Expected 1 book, got %d", len(books))
	}

	testDB.Delete(&bookToUpdate)

	var deletedBook model.Book
	testDB.First(&deletedBook, "judul = ?", "Sample Book")
	if deletedBook.ID != 0 {
		t.Error("Expected book to be deleted")
	}
}
