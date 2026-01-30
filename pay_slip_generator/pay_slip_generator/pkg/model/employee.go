package model

import (
	"fmt"
	"strconv"
	"strings"
)

// Employee represents a single row of data from the Excel file.
type Employee struct {
	// Meta
	Month string
	Year  string

	// Employee Details
	Name        string
	Designation string
	Email       string // Added Email field
	BankAcNo    string
	DOJ         string // Date of Joining
	Gender      string
	PAN         string

	// Attendance
	StandardDays string
	PayableDays  string
	LOPDays      string // Loss of Pay

	// Earnings (Rate and Amount)
	BasicPayRate   float64
	BasicPayAmount float64

	HRARate   float64
	HRAAmount float64

	OtherAllowanceRate   float64
	OtherAllowanceAmount float64

	// Deductions
	ProfessionalTax float64
	PF              float64 // Provident Fund (if any)
	IncomeTax       float64 // If any

	// Totals
	GrossEarnings   float64
	TotalDeductions float64
	NetPay          float64
}

// NetPayInWords converts the NetPay to words.
// NOTE: This is a basic implementation. For production use with varying amounts,
// a dedicated number-to-words library (like github.com/divan/num2words) is recommended.
func (e *Employee) NetPayInWords() string {
	// Basic integer part for now
	val := int(e.NetPay)
	return fmt.Sprintf("RUPEES %s ONLY", strings.ToUpper(convertNumberToWords(val)))
}

// convertNumberToWords is a simple helper for demo purposes.
// A full implementation would be much larger.
func convertNumberToWords(n int) string {
	// usage of a library would be better, but this handles simple cases or returns the number as string
	return strconv.Itoa(n)
}
