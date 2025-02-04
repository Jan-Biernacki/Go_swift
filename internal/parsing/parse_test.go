package parsing

import (
	"os"
	"strings"
	"testing"
)

type expectedRow struct {
	swiftCode   string
	isHQ        bool
	countryISO2 string
	countryName string
	bankName    string
	address     string
}

func TestParseCSV_MultipleRows(t *testing.T) {
	// CSV content with header and four rows.
	// Note: The header names don’t affect parsing since the first row is skipped.
	csvContent := `COUNTRY ISO2 CODE,SWIFT CODE,CODE TYPE,NAME,ADDRESS,TOWN NAME,COUNTRY NAME,TIME ZONE
AL,AAISALTRXXX,BIC11,UNITED BANK OF ALBANIA SH.A,"HYRJA 3 RR. DRITAN HOXHA ND. 11 TIRANA, TIRANA, 1023",TIRANA,ALBANIA,Europe/TiraneA
AW,ARUBAWAXXXX,BIC11,"ARUBA BANK, LTD",CAMACURI 12  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST,ORANJESTAD,ARUBA,America/Aruba
CL,BICHCLRMXXX,BIC11,BANCO INTERNACIONAL,  ,SANTIAGO,CHILE,Pacific/Easter
MT,TGAFMTM1001,BIC11,TGA FUNDS SICAV PLC,"  MRIEHEL, BIRKIRKARA",MRIEHEL,MALTA,Europe/Malta`

	// Create a temporary file to hold the CSV content.
	tmpFile, err := os.CreateTemp("", "swift_codes_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write the CSV content into the file.
	if _, err := tmpFile.WriteString(csvContent); err != nil {
		t.Fatalf("Failed to write CSV content: %v", err)
	}
	tmpFile.Close()

	// Call the ParseCSV function.
	parsed, err := ParseCSV(tmpFile.Name())
	if err != nil {
		t.Fatalf("ParseCSV returned error: %v", err)
	}

	// We expect four parsed rows.
	if len(parsed) != 4 {
		t.Fatalf("Expected 4 parsed rows, got %d", len(parsed))
	}

	// Define the expected results for each row.
	expectedRows := []expectedRow{
		{
			swiftCode:   "AAISALTRXXX",
			isHQ:        true,
			countryISO2: "AL",
			countryName: "ALBANIA",
			bankName:    "UNITED BANK OF ALBANIA SH.A",
			address:     "HYRJA 3 RR. DRITAN HOXHA ND. 11 TIRANA, TIRANA, 1023, TIRANA",
		},
		{
			swiftCode:   "ARUBAWAXXXX",
			isHQ:        true,
			countryISO2: "AW",
			countryName: "ARUBA",
			bankName:    "ARUBA BANK, LTD",
			address:     "CAMACURI 12  - ORANJESTAD ORANJESTAD-WEST AND ORANJESTAD-EAST, ORANJESTAD",
		},
		{
			swiftCode:   "BICHCLRMXXX",
			isHQ:        true,
			countryISO2: "CL",
			countryName: "CHILE",
			bankName:    "BANCO INTERNACIONAL",
			// The ADDRESS field is empty (after trimming) so the full address is a comma followed by the town name.
			address: ", SANTIAGO",
		},
		{
			swiftCode:   "TGAFMTM1001",
			isHQ:        false,
			countryISO2: "MT",
			countryName: "MALTA",
			bankName:    "TGA FUNDS SICAV PLC",
			// The ADDRESS field ("  MRIEHEL, BIRKIRKARA") gets trimmed to "MRIEHEL, BIRKIRKARA" then appended with the town name.
			address: "MRIEHEL, BIRKIRKARA, MRIEHEL",
		},
	}

	// Iterate over each expected row and compare with the parsed result.
	for i, expected := range expectedRows {
		actual := parsed[i]

		if actual.SwiftCode != expected.swiftCode {
			t.Errorf("Row %d: Expected SwiftCode %s, got %s", i+1, expected.swiftCode, actual.SwiftCode)
		}

		if actual.IsHeadquarter != expected.isHQ {
			t.Errorf("Row %d: Expected IsHeadquarter %v, got %v", i+1, expected.isHQ, actual.IsHeadquarter)
		}

		if actual.CountryISO2 != expected.countryISO2 {
			t.Errorf("Row %d: Expected CountryISO2 %s, got %s", i+1, expected.countryISO2, actual.CountryISO2)
		}

		if actual.CountryName != expected.countryName {
			t.Errorf("Row %d: Expected CountryName %s, got %s", i+1, expected.countryName, actual.CountryName)
		}

		if actual.BankName != expected.bankName {
			t.Errorf("Row %d: Expected BankName %s, got %s", i+1, expected.bankName, actual.BankName)
		}

		// Use TrimSpace on the address to safeguard against minor whitespace differences.
		if strings.TrimSpace(actual.Address) != expected.address {
			t.Errorf("Row %d: Expected Address %q, got %q", i+1, expected.address, actual.Address)
		}
	}
}
