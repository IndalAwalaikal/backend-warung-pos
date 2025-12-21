package main

import (
	"fmt"
	"log"

	"github.com/IndalAwalaikal/warung-pos/backend/config"
	"github.com/IndalAwalaikal/warung-pos/backend/model"
	"github.com/IndalAwalaikal/warung-pos/backend/router"
	"github.com/IndalAwalaikal/warung-pos/backend/utils"
)

func main() {
    // load env
    config.LoadEnv()

    // connect to DB (uses env vars or defaults)
    config.ConnectMySQL("")

    // auto-migrate models
    db := config.DB
    if db != nil {
        db.AutoMigrate(&model.User{}, &model.Category{}, &model.Menu{}, &model.Transaction{}, &model.TransactionItem{})
    }

    // seed admin user if not present (use env ADMIN_EMAIL / ADMIN_PASSWORD)
    if db != nil {
        adminEmail := config.GetEnv("ADMIN_EMAIL", "admin@warung.local")
        adminPass := config.GetEnv("ADMIN_PASSWORD", "admin123")
        var u model.User
        if err := db.Where("email = ?", adminEmail).First(&u).Error; err != nil {
            // create admin
            hashed, err := utils.HashPassword(adminPass)
            if err != nil {
                log.Printf("failed to hash admin password: %v", err)
            } else {
                admin := model.User{Name: "Admin", Email: adminEmail, Password: hashed, Role: "admin"}
                if err := db.Create(&admin).Error; err != nil {
                    log.Printf("failed to create admin user: %v", err)
                } else {
                    log.Printf("created admin user: %s", adminEmail)
                }
            }
        }
    }

    r := router.SetupRouter()
    port := config.GetEnv("PORT", "8080")
    fmt.Printf("starting server on :%s\n", port)
    r.Run(":" + port)
}
