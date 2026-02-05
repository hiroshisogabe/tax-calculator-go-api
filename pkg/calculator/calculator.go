package calculator

type TaxRule struct {
	State    string
	Year     int
	Category string
	Rate     float64
}

type Result struct {
	TaxAmount float64
	Total     float64
	Rate      float64
}

var mockRules = []TaxRule{
	{State: "NY", Year: 2024, Category: "electronics", Rate: 0.088},
	{State: "CA", Year: 2024, Category: "clothing", Rate: 0.075},
	{State: "TX", Year: 2024, Category: "services", Rate: 0.0},
}

func FindRate(state string, year int, category string) (float64, bool) {
	for _, rule := range mockRules {
		if rule.State == state && rule.Year == year && rule.Category == category {
			return rule.Rate, true
		}
	}
	return 0, false
}

func Calculate(amount float64, rate float64) Result {
	tax := amount * rate
	return Result{
		TaxAmount: tax,
		Total:     amount + tax,
		Rate:      rate,
	}
}
