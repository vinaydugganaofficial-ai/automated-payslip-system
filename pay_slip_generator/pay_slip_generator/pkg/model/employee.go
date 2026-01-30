package model

import (
	"fmt"
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

	UAN  string // New Field
	PFNo string // New Field - PF Account Number

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
func (e *Employee) NetPayInWords() string {
	val := int(e.NetPay)
	words := convertNumberToWords(val)
	return fmt.Sprintf("RUPEES %s ONLY", strings.ToUpper(words))
}

var (
	ones  = []string{"", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine"}
	teens = []string{"Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"}
	tens  = []string{"", "", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"}
)

func convertNumberToWords(n int) string {
	if n == 0 {
		return "Zero"
	}
	return strings.TrimSpace(recursiveConvert(n))
}

func recursiveConvert(n int) string {
	if n < 0 {
		return "Minus " + recursiveConvert(-n)
	}
	if n < 10 {
		return ones[n]
	}
	if n < 20 {
		return teens[n-10]
	}
	if n < 100 {
		return tens[n/10] + " " + ones[n%10]
	}
	if n < 1000 {
		return ones[n/100] + " Hundred " + recursiveConvert(n%100)
	}
	if n < 100000 {
		return recursiveConvert(n/1000) + " Thousand " + recursiveConvert(n%1000)
	}
	if n < 10000000 {
		return recursiveConvert(n/100000) + " Lakh " + recursiveConvert(n%100000)
	}
	return recursiveConvert(n/10000000) + " Crore " + recursiveConvert(n%10000000)
}
