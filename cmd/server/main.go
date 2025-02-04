package main

import (
	"log"

	"go_swift/internal/controllers"
	"go_swift/internal/models"
	"go_swift/internal/parsing"
	"go_swift/internal/repositories"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Init DB
	repositories.InitDB()

	// 2. Parse CSV (only once, on startup)
	parsed, err := parsing.ParseCSV("data/swift_codes.csv")
	if err != nil {
		log.Fatalf("Failed to parse CSV: %v", err)
	}
	// Insert them if the DB is empty
	seedDatabase(parsed)

	// 3. Create Gin router
	r := gin.Default()

	// 4. Define endpoints
	v1 := r.Group("/v1")
	{
		swift := v1.Group("/swift-codes")
		{
			swift.GET("/:swiftCode", controllers.GetSwiftCode)
			swift.GET("/country/:iso2", controllers.GetByCountry)
			swift.POST("", controllers.CreateSwiftCode)
			swift.DELETE("/:swiftCode", controllers.DeleteSwiftCode)
		}
	}

	// 5. Start
	r.Run(":8080")
}

// seedDatabase inserts the parsed CSV data if the DB is empty
func seedDatabase(parsed []parsing.ParsedSwift) {
	var count int64
	repositories.DB.Model(&models.SwiftCode{}).Count(&count)
	if count == 0 {
		for _, p := range parsed {
			repositories.DB.Table("swift_codes").Create(&models.SwiftCode{
				SwiftCode:     p.SwiftCode,
				BankName:      p.BankName,
				Address:       p.Address,
				CountryISO2:   p.CountryISO2,
				CountryName:   p.CountryName,
				IsHeadquarter: p.IsHeadquarter,
			})
		}
	}
}
