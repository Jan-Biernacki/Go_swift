// integration_test.go
package main

import (
	"bytes"
	"encoding/json"
	"go_swift/internal/controllers"
	"go_swift/internal/models"
	"go_swift/internal/parsing"
	"go_swift/internal/repositories"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// setupRouter replicates the routing defined in main.go.
func setupRouter() *gin.Engine {
	router := gin.Default()
	v1 := router.Group("/v1")
	{
		swift := v1.Group("/swift-codes")
		{
			swift.GET("/:swiftCode", controllers.GetSwiftCode)
			swift.GET("/country/:iso2", controllers.GetByCountry)
			swift.POST("", controllers.CreateSwiftCode)
			swift.DELETE("/:swiftCode", controllers.DeleteSwiftCode)
		}
	}
	return router
}

// seedDatabaseIfEmpty mimics the logic from main.go to seed the database if it is empty.
func seedDatabaseIfEmpty(t *testing.T) {
	var count int64
	if err := repositories.DB.Model(&models.SwiftCode{}).Count(&count).Error; err != nil {
		t.Fatalf("failed to count swift codes: %v", err)
	}
	if count == 0 {
		parsed, err := parsing.ParseCSV("data/swift_codes.csv")
		if err != nil {
			t.Fatalf("failed to parse CSV: %v", err)
		}
		for _, p := range parsed {
			repositories.DB.Create(&models.SwiftCode{
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

// TestIntegrationAPI exercises the main endpoints of the API.
func TestIntegrationAPI(t *testing.T) {
	// Load environment variables from .env file.
	if err := godotenv.Load(); err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	os.Setenv("DB_HOST", "localhost")

	// Use Gin's test mode.
	gin.SetMode(gin.TestMode)

	// Initialize the database.
	repositories.InitDB()

	repositories.DB.Exec("TRUNCATE TABLE swift_codes RESTART IDENTITY CASCADE")

	// Seed the database if it's empty.
	seedDatabaseIfEmpty(t)

	// Set up the router.
	router := setupRouter()

	// Start an HTTP test server.
	server := httptest.NewServer(router)
	defer server.Close()

	// Create an HTTP client with a timeout.
	client := &http.Client{Timeout: 10 * time.Second}

	// --- Test POST: Create a new SWIFT code entry.
	newSwift := models.SwiftCode{
		SwiftCode:     "NEWTEST33XXX",
		BankName:      "New Test Bank",
		Address:       "456 New St",
		CountryISO2:   "US",
		CountryName:   "UNITED STATES",
		IsHeadquarter: true,
	}
	newSwiftJSON, err := json.Marshal(newSwift)
	if err != nil {
		t.Fatalf("failed to marshal new swift code: %v", err)
	}
	postURL := server.URL + "/v1/swift-codes"
	resp, err := client.Post(postURL, "application/json", bytes.NewBuffer(newSwiftJSON))
	if err != nil {
		t.Fatalf("POST request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("POST: expected status 200 OK, got %d", resp.StatusCode)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	t.Logf("POST response: %s", string(body))

	// --- Test GET: Retrieve the newly created SWIFT code.
	getURL := server.URL + "/v1/swift-codes/" + newSwift.SwiftCode
	resp, err = client.Get(getURL)
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET: expected status 200 OK, got %d", resp.StatusCode)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	t.Logf("GET response: %s", string(body))

	// --- Test GET by Country: Retrieve SWIFT codes for country "US".
	getCountryURL := server.URL + "/v1/swift-codes/country/US"
	resp, err = client.Get(getCountryURL)
	if err != nil {
		t.Fatalf("GET by country request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("GET by country: expected status 200 OK, got %d", resp.StatusCode)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	t.Logf("GET by country response: %s", string(body))

	// --- Test DELETE: Delete the SWIFT code we just created.
	deleteURL := server.URL + "/v1/swift-codes/" + newSwift.SwiftCode
	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		t.Fatalf("failed to create DELETE request: %v", err)
	}
	resp, err = client.Do(req)
	if err != nil {
		t.Fatalf("DELETE request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("DELETE: expected status 200 OK, got %d", resp.StatusCode)
	}
	body, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	t.Logf("DELETE response: %s", string(body))
}
