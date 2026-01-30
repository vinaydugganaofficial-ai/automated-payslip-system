package calculator

// CalculateMonthlyIncomeTax calculates the monthly tax deduction based on FY 2025-26 slabs.
// It assumes the monthlyGross is consistent for the whole year (Simple Projection).
func CalculateMonthlyIncomeTax(monthlyGross float64) float64 {
	annualIncome := monthlyGross * 12
	annualTax := 0.0

	// Deductions/Exemptions should technically be subtracted here (e.g. Standard Deduction 50k, 80C etc.)
	// But per user request "As per Salary you should caluclate", we will stick to the basic slabs
	// on the Gross for simplicity unless implicitly asked for Standard Deduction.
	// Usually, for "New Regime" (which these slabs look like), there is a standard deduction of 75k (FY 25-26 proposed) or 50k?
	// The user gave specific slabs. I will strictly follow the USER PROVIDED SLABS on the income.

	// User Slabs:
	// Up to 2,50,000: No tax
	// 2,50,001 - 5,00,000: 5%
	// 5,00,001 - 10,00,000: 20% + 12,500
	// Above 10,00,000: 30% + 1,12,500

	if annualIncome <= 250000 {
		annualTax = 0
	} else if annualIncome <= 500000 {
		annualTax = (annualIncome - 250000) * 0.05
	} else if annualIncome <= 1000000 {
		// 5% of (5L - 2.5L) = 12,500. This matches the user's fixed amount.
		annualTax = 12500 + (annualIncome-500000)*0.20
	} else {
		// 12,500 + 20% of (10L - 5L) = 12,500 + 1,00,000 = 1,12,500. Matches user's fixed amount.
		annualTax = 112500 + (annualIncome-1000000)*0.30
	}

	return annualTax / 12
}
