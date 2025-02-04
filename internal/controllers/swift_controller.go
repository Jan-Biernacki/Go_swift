package controllers

import (
	"net/http"
	"strings"

	"go_swift/internal/models"
	repo "go_swift/internal/repositories"

	"github.com/gin-gonic/gin"
)

// GET /v1/swift-codes/:swiftCode
func GetSwiftCode(c *gin.Context) {
	code := c.Param("swiftCode")

	var result models.SwiftCode
	if err := repo.DB.Where("swift_code = ?", code).First(&result).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Swift code not found"})
		return
	}

	// If it's a HQ, find branches
	if result.IsHeadquarter {
		prefix := result.SwiftCode[:8]
		var branches []models.SwiftCode
		repo.DB.Where("swift_code LIKE ? AND id <> ?", prefix+"%", result.ID).Find(&branches)

		c.JSON(http.StatusOK, gin.H{
			"address":       result.Address,
			"bankName":      result.BankName,
			"countryISO2":   result.CountryISO2,
			"countryName":   result.CountryName,
			"isHeadquarter": result.IsHeadquarter,
			"swiftCode":     result.SwiftCode,
			"branches":      branches,
		})
	} else {
		// If it's a branch, just return it
		c.JSON(http.StatusOK, result)
	}
}

// GET /v1/swift-codes/country/:iso2
func GetByCountry(c *gin.Context) {
	iso2 := strings.ToUpper(c.Param("iso2"))
	var codes []models.SwiftCode
	repo.DB.Where("country_iso2 = ?", iso2).Find(&codes)
	if len(codes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data for that country"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"countryISO2": iso2,
		"countryName": codes[0].CountryName,
		"swiftCodes":  codes,
	})
}

// POST /v1/swift-codes
func CreateSwiftCode(c *gin.Context) {
	var input models.SwiftCode
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	// Insert
	if err := repo.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create swift code"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Created"})
}

// DELETE /v1/swift-codes/:swiftCode
func DeleteSwiftCode(c *gin.Context) {
	code := c.Param("swiftCode")
	bankName := c.Query("bankName")
	iso2 := strings.ToUpper(c.Query("countryISO2"))

	var sc models.SwiftCode
	// Check if record exists
	err := repo.DB.Where("swift_code = ? AND bank_name = ? AND country_iso2 = ?", code, bankName, iso2).First(&sc).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching record found"})
		return
	}
	repo.DB.Delete(&sc)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
