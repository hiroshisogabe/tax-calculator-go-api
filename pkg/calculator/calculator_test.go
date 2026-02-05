package calculator

import "testing"

func TestFindRate(t *testing.T) {
	tests := []struct {
		name          string
		state         string
		year          int
		category      string
		expectedRate  float64
		expectedFound bool
	}{
		{"Valid NY Rule", "NY", 2024, "electronics", 0.088, true},
		{"Valid CA Rule", "CA", 2024, "clothing", 0.075, true},
		{"Invalid State", "ZZ", 2024, "electronics", 0, false},
		{"Wrong Year", "NY", 1999, "electronics", 0, false},
		{"Empty Category", "NY", 2024, "", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rate, found := FindRate(tt.state, tt.year, tt.category)

			if found != tt.expectedFound {
				t.Fatalf("FindRate() found = %v, want %v", found, tt.expectedFound)
			}

			if rate != tt.expectedRate {
				t.Errorf("FindRate() rate = %v, want %v", rate, tt.expectedRate)
			}
		})
	}
}

func TestCalculate(t *testing.T) {
	tests := []struct {
		name          string
		amount        float64
		rate          float64
		expectedTax   float64
		expectedTotal float64
	}{
		{"Standard 10% tax", 100, 0.10, 10, 110},
		{"Zero tax", 100, 0.0, 0, 100},
		{"High tax", 200, 0.25, 50, 250},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Calculate(tt.amount, tt.rate)

			if got.TaxAmount != tt.expectedTax {
				t.Errorf("Calculate(%v, %v) TaxAmount = %v; want %v", tt.amount, tt.rate, got.TaxAmount, tt.expectedTax)
			}

			if got.Total != tt.expectedTotal {
				t.Errorf("Calculate(%v, %v) Total = %v; want %v", tt.amount, tt.rate, got.Total, tt.expectedTotal)
			}
		})
	}
}
