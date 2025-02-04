package models

import (
	"gorm.io/gorm"
)

// SwiftCode represents the structure of a parsed data record.
type SwiftCode struct {
	gorm.Model
	SwiftCode     string `gorm:"column:swift_code;uniqueIndex"`
	BankName      string `gorm:"column:bank_name"`
	Address       string `gorm:"column:address"`
	CountryISO2   string `gorm:"column:country_iso2"`
	CountryName   string `gorm:"column:country_name"`
	IsHeadquarter bool   `gorm:"column:is_headquarter"`
}

func (SwiftCode) TableName() string {
	return "swift_codes"
}
