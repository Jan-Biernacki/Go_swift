package controllers

import (
	"fmt"
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

	// If HQ, find branches
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
		// If branch, just return it
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
	// Retrieve the swift code from the URL parameter.
	code := c.Param("swiftCode")
	fmt.Printf("Attempting to delete record with swift code: '%s'\n", code)

	var sc models.SwiftCode
	// Look up the record solely based on the swift code.
	if err := repo.DB.Where("swift_code = ?", code).First(&sc).Error; err != nil {
		// Log the error for debugging.
		fmt.Printf("Record not found. Error: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "No matching record found"})
		return
	}

	// Attempt deletion and check for errors.
	if err := repo.DB.Delete(&sc).Error; err != nil {
		fmt.Printf("Error deleting record: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete swift code"})
		return
	}

	fmt.Printf("Record with swift code '%s' deleted successfully.\n", code)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}
