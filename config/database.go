package config

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectMySQL connects to MySQL using env vars and stores the DB handle in package variable
func ConnectMySQL(dsnOverride string) {
    var dsn string
    if dsnOverride != "" {
        dsn = dsnOverride
    } else {
        user := GetEnv("DB_USER", "root")
        pass := GetEnv("DB_PASS", "")
        host := GetEnv("DB_HOST", "127.0.0.1")
        port := GetEnv("DB_PORT", "3306")
        name := GetEnv("DB_NAME", "warung_pos")
        dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)
    }

    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect database: %v", err)
    }
    DB = db
}
