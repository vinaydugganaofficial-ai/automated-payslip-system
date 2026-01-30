package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()
	sheetName := "Sheet1"

	// Create headers
	headers := []string{
		"Month", "Year", "Emp Name", "Email", "Designation", "Bank Ac No", "DOJ", "Gender", "PAN",
		"Standard Days", "Payable Days", "LOP Days",
		"Basic Pay Rate", "HRA Rate", "Other Allowance Rate",
		"Basic Pay", "HRA", "Other Allowance",
		"Professional Tax", "PF", "Income Tax",
		"Gross Earnings", "Total Deductions", "Net Pay",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, h)
	}

	// Sample Data
	data := []interface{}{
		"Dec", "2024", "Vinay", "vinayopbr@gmail.com", "Software Engineer", "1234567890", "2023-01-01", "Male", "ABCDE1234F",
		"31", "31", "0",
		"50000", "20000", "10000",
		"50000", "20000", "10000",
		"200", "1800", "1000",
		"80000", "3000", "77000",
	}

	for i, v := range data {
		cell, _ := excelize.CoordinatesToCellName(i+1, 2)
		f.SetCellValue(sheetName, cell, v)
	}

	if err := f.SaveAs("employees.xlsx"); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Created employees.xlsx with sample data.")
}
