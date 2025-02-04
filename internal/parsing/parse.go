package parsing

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

// RawSwift holds columns as read from CSV.
type RawSwift struct {
	CountryISO2 string
	SwiftCode   string
	CodeType    string
	Name        string
	Address     string
	TownName    string
	CountryName string
	TimeZone    string
}

// ParsedSwift is the final shape we want after parsing.
type ParsedSwift struct {
	SwiftCode     string
	BankName      string
	Address       string
	CountryISO2   string
	CountryName   string
	IsHeadquarter bool
}

// ParseCSV reads rows from the CSV and returns a slice of ParsedSwift.
func ParseCSV(filePath string) ([]ParsedSwift, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	// Skip header row
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	var results []ParsedSwift

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Trim & uppercase relevant fields
		countryISO2 := strings.ToUpper(strings.TrimSpace(row[0]))
		swiftCode := strings.TrimSpace(row[1])
		// Validate swift code length 8 is in case that the user wants to see HQ as well as affilieaded branches
		if len(swiftCode) != 8 && len(swiftCode) != 11 {
			return nil, fmt.Errorf("invalid swift code '%s': must be either 8 or 11 characters", swiftCode)
		}
		// row[2] = codeType
		name := strings.TrimSpace(row[3])
		address := strings.TrimSpace(row[4])
		townName := strings.TrimSpace(row[5])
		countryName := strings.ToUpper(strings.TrimSpace(row[6]))
		// row[7] = timeZone

		// Decide if it's HQ
		isHQ := false
		if len(swiftCode) == 8 || (len(swiftCode) == 11 && strings.HasSuffix(swiftCode, "XXX")) {
			isHQ = true
		}

		// Combine address + town if you like:
		fullAddress := address
		if townName != "" {
			fullAddress += ", " + townName
		}

		results = append(results, ParsedSwift{
			SwiftCode:     swiftCode,
			BankName:      name,
			Address:       fullAddress,
			CountryISO2:   countryISO2,
			CountryName:   countryName,
			IsHeadquarter: isHQ,
		})
	}

	return results, nil
}
