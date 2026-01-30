package generator

import (
	"fmt"
	"pay_slip_generator/pkg/model"

	"github.com/jung-kurt/gofpdf"
)

// GeneratePaySlip creates a PDF pay slip for the given employee.
func GeneratePaySlip(emp model.Employee, outputDir string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// --- Styles ---
	pdf.SetFont("Arial", "", 10) // Basic font

	// --- Header ---
	// Logo
	// Using the provided logo image "logo.png"
	// Adjust coordinates and width (40mm) as needed to match the look
	pdf.Image("logo.png", 15, 10, 40, 0, false, "", 0, "")

	// Reset Color (just in case)
	pdf.SetTextColor(0, 0, 0)

	// Company Address (Right Aligned)
	// Align closer to right margin (A4 width 210, margin 10/15 -> ~195)
	pdf.SetXY(110, 15)
	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(85, 5, "AbegaTech Pvt. Ltd.")
	pdf.Ln(5)
	pdf.SetX(110)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(85, 4, "P No 147, Floor 1 Rd No7, Sri Madhavam,\nMadeenaguda, Miyapur, Hyderabad 500049", "", "L", false)

	pdf.SetY(40) // Space before content

	// --- Title Box ---
	// Grey background title
	pdf.SetFillColor(230, 230, 230)
	pdf.SetDrawColor(0, 0, 0) // Black borders
	pdf.SetFont("Arial", "", 10)
	pdf.SetLineWidth(0.3)

	// Payslip for : Month Year (Right Aligned in the box)
	// Detailed Grid border box begins
	pdf.SetX(10)
	pdf.CellFormat(190, 7, fmt.Sprintf("Payslip for : %s %s", emp.Month, emp.Year), "1", 1, "R", true, 0, "")

	// --- Employee Details Grid ---
	pdf.SetFont("Arial", "", 9)
	h := 7.0 // Row height

	// Widths: Label 25, Value 65, Label 30, Value 70 = 190 total

	// Row 1
	pdf.SetX(10)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(25, h, " Emp Name", "L", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(65, h, emp.Name, "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(30, h, "DOJ", "", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(70, h, "  "+emp.DOJ, "R", 1, "L", false, 0, "")

	// Row 2
	pdf.SetX(10)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(25, h, " Designation", "L", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(65, h, emp.Designation, "", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(30, h, "Gender", "", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(70, h, "  "+emp.Gender, "R", 1, "L", false, 0, "")

	// Row 3 (Bottom border for first block)
	pdf.SetX(10)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(25, h, " Bank Ac. No.", "LB", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(65, h, emp.BankAcNo, "B", 0, "L", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(30, h, "PAN", "B", 0, "R", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(70, h, "  "+emp.PAN, "RB", 1, "L", false, 0, "")

	// --- Attendance Info ---
	// Grey background
	pdf.SetFillColor(230, 230, 230)
	pdf.SetFont("Arial", "", 9)
	pdf.SetX(10)
	attendanceText := fmt.Sprintf("Standard Days: %s          Payable days: %s          Loss of Pay Days : %s", emp.StandardDays, emp.PayableDays, emp.LOPDays)
	pdf.CellFormat(190, 7, attendanceText, "1", 1, "L", true, 0, "")

	// --- Earnings & Deductions Tables ---
	pdf.SetFont("Arial", "B", 9)

	// Header
	pdf.SetX(10)
	pdf.CellFormat(55, 8, " Earnings", "LBT", 0, "L", false, 0, "")

	// Stacked Header trick using XY reset or just simple line
	// "Standard Rate" logic
	origX, origY := pdf.GetX(), pdf.GetY()
	pdf.CellFormat(20, 8, "", "BT", 0, "C", false, 0, "") // Frame
	pdf.SetXY(origX, origY)
	pdf.SetFont("Arial", "B", 8)
	pdf.CellFormat(20, 4, "Standard", "", 0, "C", false, 0, "") // Top
	pdf.SetXY(origX, origY+4)
	pdf.CellFormat(20, 4, "Rate", "", 0, "C", false, 0, "") // Bottom
	pdf.SetXY(origX+20, origY)                              // Next col
	pdf.SetFont("Arial", "B", 9)

	pdf.CellFormat(20, 8, "Amount", "BTR", 0, "R", false, 0, "")

	// Deductions Header
	pdf.CellFormat(55, 8, " Deductions", "BT", 0, "L", false, 0, "")
	pdf.CellFormat(40, 8, "Total", "BTR", 1, "R", false, 0, "")

	// --- Table Content ---
	pdf.SetFont("Arial", "", 9)

	// Helper for rows
	drawRow := func(earnLabel string, earnRate, earnAmt float64, dedLabel string, dedAmt float64) {
		pdf.SetX(10)

		// Earnings
		pdf.CellFormat(55, 6, " "+earnLabel, "L", 0, "L", false, 0, "")

		rateStr := ""
		if earnRate > 0 {
			rateStr = fmt.Sprintf("%.2f", earnRate)
		}
		pdf.CellFormat(20, 6, rateStr, "", 0, "R", false, 0, "")

		amtStr := ""
		if earnAmt > 0 || earnRate > 0 {
			amtStr = fmt.Sprintf("%.2f", earnAmt)
		}
		pdf.CellFormat(20, 6, amtStr, "R", 0, "R", false, 0, "")

		// Deductions
		pdf.CellFormat(55, 6, " "+dedLabel, "", 0, "L", false, 0, "")

		dedStr := ""
		if dedAmt > 0 {
			dedStr = fmt.Sprintf("%.2f", dedAmt)
		}
		pdf.CellFormat(40, 6, dedStr, "R", 1, "R", false, 0, "")
	}

	// Rows
	drawRow("Basic Pay", emp.BasicPayRate, emp.BasicPayAmount, "Professional Tax", emp.ProfessionalTax)
	drawRow("House Rent Allowance", emp.HRARate, emp.HRAAmount, "Provident Fund", emp.PF)
	drawRow("Other Allowance", emp.OtherAllowanceRate, emp.OtherAllowanceAmount, "Income Tax", emp.IncomeTax)
	// Empty rows
	drawRow("", 0, 0, "", 0)
	drawRow("", 0, 0, "", 0)

	// --- Totals ---
	pdf.SetFont("Arial", "B", 9)
	pdf.SetX(10)
	// Gross Earnings
	pdf.CellFormat(55, 8, " Gross Earnings", "LTB", 0, "L", false, 0, "")
	pdf.CellFormat(20, 8, fmt.Sprintf("%.2f", emp.GrossEarnings), "TB", 0, "R", false, 0, "")
	pdf.CellFormat(20, 8, fmt.Sprintf("%.2f", emp.GrossEarnings), "TBR", 0, "R", false, 0, "")

	// Total Deductions
	pdf.CellFormat(55, 8, " Total Deductions", "TB", 0, "L", false, 0, "")
	pdf.CellFormat(40, 8, fmt.Sprintf("%.0f", emp.TotalDeductions), "TBR", 1, "R", false, 0, "")

	// --- Net Pay ---
	pdf.SetX(10)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(30, 10, " NET PAY", "LTB", 0, "L", false, 0, "")
	pdf.CellFormat(25, 10, fmt.Sprintf("%.2f", emp.NetPay), "TB", 0, "L", false, 0, "")

	// Words
	pdf.SetFont("Arial", "", 8)
	words := emp.NetPayInWords()
	pdf.CellFormat(135, 10, "("+words+")", "TBR", 1, "L", false, 0, "")

	pdf.Ln(10)

	// --- Footer ---
	pdf.SetX(10)
	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(190, 5, "** This is computer generated payslip and doesn't require signature and stamp", "", 1, "C", false, 0, "")

	// --- Download/Print Button ---
	// Visual button
	pdf.SetY(pdf.GetY() + 5)
	pdf.SetX(90)
	pdf.SetFillColor(0, 150, 150) // Teal button
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 10)

	// Draw Button
	pdf.CellFormat(30, 10, "PRINT", "1", 0, "C", true, 0, "")

	// Add Link covering the button area
	linkID := pdf.AddLink()
	pdf.Link(90, pdf.GetY(), 30, 10, linkID)

	// Javascript to trigger print dialog
	pdf.SetJavascript("function Print() { print(); }")
	// Note: Link action to JS is not auto-wired here without advanced usage,
	// but the JS is embedded in the PDF so Opening it might trigger or CRTL+P fits best.
	// Ideally user clicks Print icon in viewer.
	// Visual cue only for now as requested.

	// Write file
	outfile := fmt.Sprintf("%s/%s_%s_%s.pdf", outputDir, emp.Name, emp.Month, emp.Year)
	return pdf.OutputFileAndClose(outfile)
}
