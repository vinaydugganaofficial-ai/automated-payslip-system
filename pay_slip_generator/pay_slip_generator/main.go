package main

import (
	"fmt"
	"log"
	"os"
	"pay_slip_generator/pkg/generator"
	"pay_slip_generator/pkg/reader"
)

func main() {
	fmt.Println("Pay Slip Generator started...")

	// Configuration
	inputFile := "employee_payslip_data_10_employees.xlsx" // Updated filename
	outputDir := "output"

	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		fmt.Printf("Warning: Input file '%s' not found.\n", inputFile)
		fmt.Println("Please ensure 'employees.xlsx' is in the project folder:")
		curr, _ := os.Getwd()
		fmt.Println(curr)
		// Ensure we don't crash if checking user project
		// return
	}

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Could not create output directory: %v", err)
	}

	// 1. Read Employees
	fmt.Printf("Reading data from %s...\n", inputFile)
	employees, err := reader.ReadEmployees(inputFile)
	if err != nil {
		// If file not found or invalid, log and exit
		log.Printf("Error reading Excel file: %v. \n(Note: If file is missing, paste 'employees.xlsx' in the folder)", err)
		return
	}
	fmt.Printf("Found %d employee records.\n", len(employees))

	// 2. Generate PDFs
	for _, emp := range employees {
		err := generator.GeneratePaySlip(emp, outputDir)
		if err != nil {
			log.Printf("Failed to generate PDF for %s: %v", emp.Name, err)
		} else {
			fmt.Printf("Generated: %s_%s_%s.pdf\n", emp.Name, emp.Month, emp.Year)
		}
	}

	fmt.Println("Processing complete. Check 'output' directory.")
}
