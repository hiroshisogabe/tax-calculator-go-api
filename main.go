package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// --- 1. Data Structures ---

// TaxRequest matches your Zod schema inputs
type TaxRequest struct {
	Amount          float64 `json:"amount"`
	State           string  `json:"state"`
	Year            int     `json:"year"`
	ProductCategory string  `json:"productCategory"`
}

// TaxResult matches the successful calculation data
type TaxResult struct {
	BaseAmount float64 `json:"baseAmount"`
	TaxAmount  float64 `json:"taxAmount"`
	Total      float64 `json:"total"`
	Rate       float64 `json:"rate"`
	State      string  `json:"state"`
	Year       int     `json:"year"`
}

// APIResponse matches your ActionResponse structure partially
type APIResponse struct {
	Success bool       `json:"success"`
	Data    *TaxResult `json:"data,omitempty"` // pointer so it can be null
	Error   string     `json:"error,omitempty"`
}

// --- 2. Mock Database & Logic ---

// TaxRule represents a row in your rules database
type TaxRule struct {
	State    string
	Year     int
	Category string
	Rate     float64
}

// mockRules simulates your database/service lookup
var mockRules = []TaxRule{
	{State: "NY", Year: 2024, Category: "electronics", Rate: 0.088},
	{State: "CA", Year: 2024, Category: "clothing", Rate: 0.075},
	{State: "TX", Year: 2024, Category: "services", Rate: 0.0},
}

// findRate simulates your findTax service
func findRate(state string, year int, category string) (float64, bool) {
	for _, rule := range mockRules {
		if rule.State == state && rule.Year == year && rule.Category == category {
			return rule.Rate, true
		}
	}
	return 0, false
}

// --- 3. HTTP Handler ---

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	// CORS setup
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode JSON
	var req TaxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid JSON format")
		return
	}

	// --- Validation ---
	req.State = strings.TrimSpace(strings.ToUpper(req.State)) // Normalize state

	if req.Amount <= 0 {
		sendError(w, "Amount must be greater than zero")
		return
	}
	if len(req.State) < 2 {
		sendError(w, "State code is required (e.g., NY)")
		return
	}
	if req.Year < 1000 || req.Year > 9999 {
		sendError(w, "Year must be a 4-digit number")
		return
	}
	if req.ProductCategory == "" {
		sendError(w, "Category is required")
		return
	}

	// --- Business Logic ---
	rate, found := findRate(req.State, req.Year, req.ProductCategory)

	if !found {
		// Mimicking your "Tax rules not available" error
		msg := fmt.Sprintf("Tax rules for %s in %d are not available for the %s category.", req.State, req.Year, req.ProductCategory)
		sendError(w, msg)
		return
	}

	// Calculation
	taxAmount := req.Amount * rate
	total := req.Amount + taxAmount

	// Success Response
	resp := APIResponse{
		Success: true,
		Data: &TaxResult{
			BaseAmount: req.Amount,
			TaxAmount:  taxAmount,
			Total:      total,
			Rate:       rate,
			State:      req.State,
			Year:       req.Year,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Helper to send JSON errors easily
func sendError(w http.ResponseWriter, message string) {
	w.Header().Set("Content-Type", "application/json")
	// You might want to use 400 Bad Request, but to match your ActionResponse structure (success: false), we can keep 200 or use 400.
	// Standard REST APIs usually return 400 for validation errors.
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(APIResponse{
		Success: false,
		Error:   message,
	})
}

func main() {
	http.HandleFunc("/calculate", calculateHandler)

	port := ":8080"
	fmt.Printf("Tax Engine running on http://localhost%s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
