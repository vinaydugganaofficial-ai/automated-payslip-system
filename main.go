package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"pay_slip_generator/pkg/calculator"
	"pay_slip_generator/pkg/generator"
	"pay_slip_generator/pkg/model"
	"pay_slip_generator/pkg/reader"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func main() {
	// 0. Parse Flags
	inputFlag := flag.String("input", "", "Path to input file (CSV or Excel)")
	dryRunFlag := flag.Bool("dry-run", false, "Run without sending emails")
	flag.Parse()

	// 1. Load Configuration
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	senderEmail := os.Getenv("SMTP_EMAIL")
	senderPassword := os.Getenv("SMTP_PASSWORD")
	fromName := os.Getenv("SMTP_FROM_NAME")
	if fromName == "" {
		fromName = "HR Team"
	}

	// In production, don't silently default to Gmail — force correct config.
	if !*dryRunFlag {
		if smtpHost == "" || smtpPortStr == "" || senderEmail == "" || senderPassword == "" {
			log.Fatal("SMTP_HOST, SMTP_PORT, SMTP_EMAIL, SMTP_PASSWORD must be set in .env")
		}
	}

	smtpPort := 587
	if smtpPortStr != "" {
		p, err := strconv.Atoi(smtpPortStr)
		if err != nil {
			log.Fatalf("Invalid SMTP_PORT=%q: %v", smtpPortStr, err)
		}
		smtpPort = p
	}

	// 2. Setup Directories
	inputFile := "employee_payslip_data_10_employees"
	if *inputFlag != "" {
		inputFile = *inputFlag
	} else if _, err := os.Stat("employees.csv"); err == nil {
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
	var d *gomail.Dialer
	var s gomail.SendCloser

	if !*dryRunFlag {
		d = gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

		// GoDaddy SMTP commonly uses implicit SSL on 465
		if smtpPort == 465 {
			d.SSL = true
		}

		// Helps TLS handshake on many servers
		d.TLSConfig = &tls.Config{
			ServerName: smtpHost,
		}

		s, err = d.Dial() // IMPORTANT: "=" not ":="
		if err != nil {
			log.Fatalf("Failed to connect to SMTP server %s:%d: %v", smtpHost, smtpPort, err)
		}
		defer s.Close()
	} else {
		fmt.Println("[DRY RUN] Skipping SMTP connection.")
	}

	// 5. Process Each Employee
	for _, emp := range employees {
		// --- Business Logic ---
		if emp.LOPDays == "" {
			emp.LOPDays = "0"
		}

		emp.IncomeTax = calculator.CalculateMonthlyIncomeTax(emp.GrossEarnings)
		emp.TotalDeductions = emp.ProfessionalTax + emp.PF + emp.IncomeTax
		emp.NetPay = emp.GrossEarnings - emp.TotalDeductions
		// ----------------------

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
		if *dryRunFlag {
			fmt.Printf("  [DRY RUN] Email sending skipped for %s (%s)\n", emp.Name, emp.Email)
			continue
		}

		m := gomail.NewMessage()
		m.SetAddressHeader("From", senderEmail, fromName)
		m.SetHeader("To", emp.Email)
		m.SetHeader("Subject", fmt.Sprintf("Payslip for %s %s", emp.Month, emp.Year))
		m.SetBody("text/plain",
			fmt.Sprintf("Dear %s,\n\nPlease find attached your payslip for %s %s.\n\nBest Regards,\n%s",
				emp.Name, emp.Month, emp.Year, fromName),
		)
		m.Attach(pdfPath)

		// Send using the persistent connection
		if err := gomail.Send(s, m); err != nil {
			log.Printf("  [ERROR] Failed to send email to %s: %v\n", emp.Email, err)

			// Simple retry logic: Re-dial if connection dropped
			_ = s.Close()
			if s, err = d.Dial(); err == nil {
				if errRetry := gomail.Send(s, m); errRetry != nil {
					log.Printf("  [ERROR] Retry failed for %s: %v\n", emp.Email, errRetry)
				} else {
					fmt.Printf("  [SUCCESS] Email sent to %s (Retry)\n", emp.Email)
				}
			} else {
				log.Printf("  [ERROR] Reconnect failed: %v\n", err)
			}
		} else {
			fmt.Printf("  [SUCCESS] Email sent to %s\n", emp.Email)
		}
	}

	fmt.Println("All tasks completed.")
}
