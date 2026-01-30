package reader

import (
	"fmt"
	"log"
	"pay_slip_generator/pkg/model"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ReadEmployees reads the excel file and returns a list of employees.
func ReadEmployees(filePath string) ([]model.Employee, error) {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Get all rows in the first sheet.
	sheetName := f.GetSheetName(0)
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("excel file is empty or missing header")
	}

	// Map headers to indices with normalizer
	headerMap := make(map[string]int)
	for i, cell := range rows[0] {
		normalized := strings.ToLower(strings.TrimSpace(cell))
		headerMap[normalized] = i
	}

	// Internal helper to get value safely
	getVal := func(row []string, possibleNames ...string) string {
		for _, name := range possibleNames {
			idx, ok := headerMap[strings.ToLower(name)]
			if ok && idx < len(row) {
				return row[idx]
			}
		}
		return ""
	}

	getFloat := func(row []string, possibleNames ...string) float64 {
		valStr := getVal(row, possibleNames...)
		valStr = strings.ReplaceAll(valStr, ",", "") // Remove commas
		val, _ := strconv.ParseFloat(valStr, 64)
		return val
	}

	var employees []model.Employee

	for _, row := range rows[1:] {
		if len(row) == 0 {
			continue
		}

		// Basic validation - if Name is empty, skip
		name := getVal(row, "Emp Name", "Name", "Employee Name", "Employee")
		if name == "" {
			// Try looking for ID if name is missing? Or just skip
			continue
		}

		// We assume headers based on the screenshot fields
		emp := model.Employee{
			Month:       getVal(row, "Month"),
			Year:        getVal(row, "Year"),
			Name:        name,
			Designation: getVal(row, "Designation", "Role", "Position"),
			Email:       getVal(row, "Email", "Email Address", "E-mail"),
			BankAcNo:    getVal(row, "Bank Ac No", "Bank Account", "Account No"),
			DOJ:         getVal(row, "DOJ", "Date of Joining", "Joining Date"),
			Gender:      getVal(row, "Gender", "Sex"),
			PAN:         getVal(row, "PAN", "PAN Number"),

			StandardDays: getVal(row, "Standard Days", "Std Days", "Total Days"),
			PayableDays:  getVal(row, "Payable Days", "Paid Days"),
			LOPDays:      getVal(row, "Loss of Pay Days", "LOP", "Absent"),

			// Earnings
			BasicPayRate:         getFloat(row, "Basic Pay Rate", "Basic Rate"),
			BasicPayAmount:       getFloat(row, "Basic Pay", "Basic Pay Amount", "Basic"),
			HRARate:              getFloat(row, "HRA Rate"),
			HRAAmount:            getFloat(row, "HRA", "House Rent Allowance"),
			OtherAllowanceRate:   getFloat(row, "Other Allowance Rate", "Other Allw Rate"),
			OtherAllowanceAmount: getFloat(row, "Other Allowance", "Other Allowance Amount", "Other Allw"),

			// Deductions
			ProfessionalTax: getFloat(row, "Professional Tax", "Prof Tax", "PT"),
			PF:              getFloat(row, "PF", "Provident Fund"),
			IncomeTax:       getFloat(row, "Income Tax", "IT", "TDS"),

			// Totals
			GrossEarnings:   getFloat(row, "Gross Earnings", "Gross Pay", "Total Earnings"),
			TotalDeductions: getFloat(row, "Total Deductions", "Total Ded"),
			NetPay:          getFloat(row, "Net Pay", "Net Salary"),
		}

		// Auto-calculate totals if missing (robustness)
		if emp.GrossEarnings == 0 {
			emp.GrossEarnings = emp.BasicPayAmount + emp.HRAAmount + emp.OtherAllowanceAmount
		}
		if emp.TotalDeductions == 0 {
			emp.TotalDeductions = emp.ProfessionalTax + emp.PF + emp.IncomeTax
		}
		if emp.NetPay == 0 {
			emp.NetPay = emp.GrossEarnings - emp.TotalDeductions
		}

		// Defaults if Month/Year missing in row (maybe take from filename or user input later? For now hardcode or leave empty)
		if emp.Month == "" {
			emp.Month = "Dec"
		} // Fallback for testing
		if emp.Year == "" {
			emp.Year = "2024"
		}

		employees = append(employees, emp)
	}

	// Debug log if no employees found but headers existed
	if len(employees) == 0 {
		log.Println("Warning: No valid employee rows found. Check column headers.")
		log.Printf("Found headers: %v", headerMap)
	}

	return employees, nil
}
