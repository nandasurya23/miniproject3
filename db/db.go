// db/db.go
package db

import (
    "fmt"

    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func ConnectDB(host, port, user, password, dbname string) (*gorm.DB, error) {
    dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname) // Ubah ini
    return gorm.Open(mysql.Open(dsn), &gorm.Config{}) 
}
