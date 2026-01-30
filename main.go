package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"pay_slip_generator/pkg/generator"
	"pay_slip_generator/pkg/model"
	"pay_slip_generator/pkg/reader"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func main() {
	// 1. Load Configuration
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	smtpHost := "smtp.gmail.com"
	smtpPort := 587
	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPassword := os.Getenv("SMTP_PASSWORD")

	if senderEmail == "" || senderPassword == "" {
		log.Fatal("Error: SMTP_EMAIL and SMTP_PASSWORD must be set in .env")
	}

	// 2. Setup Directories
	// Check for CSV first, default to Excel
	inputFile := "employee_payslip_data_10_employees"
	if _, err := os.Stat("employees.csv"); err == nil {
		inputFile = "employees.csv"
	}

	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Could not create output directory: %v", err)
	}

	// 3. Read Employees
	fmt.Printf("Reading employees from %s...\n", inputFile)

	var employees []model.Employee
	var err error

	// Check extension
	if filepath.Ext(inputFile) == ".csv" {
		employees, err = reader.ReadEmployeesFromCSV(inputFile)
		if err != nil {
			log.Fatalf("Error reading CSV: %v", err)
		}
	} else {
		employees, err = reader.ReadEmployees(inputFile)
		if err != nil {
			log.Fatalf("Error reading Excel: %v", err)
		}
	}
	fmt.Printf("Found %d employees.\n", len(employees))

	// 4. Setup SMTP Connection (Persistent)
	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	s, err := d.Dial()
	if err != nil {
		log.Fatalf("Failed to connect to SMTP server: %v", err)
	}
	defer s.Close()

	// 5. Process Each Employee
	for _, emp := range employees {
		// A. Generate PDF
		fmt.Printf("Processing %s (%s)...\n", emp.Name, emp.Email)
		err := generator.GeneratePaySlip(emp, outputDir)
		if err != nil {
			log.Printf("  [ERROR] Failed to generate PDF for %s: %v\n", emp.Name, err)
			continue
		}

		pdfPath := filepath.Join(outputDir, fmt.Sprintf("%s_%s_%s.pdf", emp.Name, emp.Month, emp.Year))

		// B. Send Email
		if emp.Email == "" {
			log.Printf("  [SKIP] No email address for %s\n", emp.Name)
			continue
		}

		m := gomail.NewMessage()
		m.SetHeader("From", senderEmail)
		m.SetHeader("To", emp.Email)
		m.SetHeader("Subject", fmt.Sprintf("Payslip for %s %s", emp.Month, emp.Year))
		m.SetBody("text/plain", fmt.Sprintf("Dear %s,\n\nPlease find attached your payslip for %s %s.\n\nBest Regards,\nHR Team", emp.Name, emp.Month, emp.Year))
		m.Attach(pdfPath)

		// Send using the persistent connection
		if err := gomail.Send(s, m); err != nil {
			log.Printf("  [ERROR] Failed to send email to %s: %v\n", emp.Email, err)
			// Simple retry logic: Re-dial if connection dropped
			s.Close()
			if s, err = d.Dial(); err == nil {
				if errRetry := gomail.Send(s, m); errRetry != nil {
					log.Printf("  [ERROR] Retry failed for %s: %v\n", emp.Email, errRetry)
				} else {
					fmt.Printf("  [SUCCESS] Email sent to %s (Retry)\n", emp.Email)
				}
			}
		} else {
			fmt.Printf("  [SUCCESS] Email sent to %s\n", emp.Email)
		}
	}

	fmt.Println("All tasks completed.")
}
