package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hiroshisogabe/tax-calculator-go-api/pkg/calculator"
)

type TaxRequest struct {
	Amount          float64 `json:"amount"`
	State           string  `json:"state"`
	Year            int     `json:"year"`
	ProductCategory string  `json:"productCategory"`
}

func (req *TaxRequest) Validate() error {
	req.State = strings.TrimSpace(strings.ToUpper(req.State))

	if req.Amount <= 0 {
		return fmt.Errorf("Amount must be greater than zero")
	}
	if len(req.State) < 2 {
		return fmt.Errorf("State code is required (e.g., NY)")
	}
	if req.Year < 1000 || req.Year > 9999 {
		return fmt.Errorf("Year must be a 4-digit number")
	}
	if req.ProductCategory == "" {
		return fmt.Errorf("Category is required")
	}

	return nil
}

type TaxResponse struct {
	BaseAmount float64 `json:"baseAmount"`
	TaxAmount  float64 `json:"taxAmount"`
	Total      float64 `json:"total"`
	Rate       float64 `json:"rate"`
	State      string  `json:"state"`
	Year       int     `json:"year"`
}

type APIResponse struct {
	Success bool         `json:"success"`
	Data    *TaxResponse `json:"data,omitempty"`  // pointer to explicitly represent null when no result handled by omitempty
	Error   string       `json:"error,omitempty"` // omitempty handles empty string (no error case)
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	setupCORS(&w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// --- Decode JSON ---
	var req TaxRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, "Invalid JSON format")
		return
	}

	if err := req.Validate(); err != nil {
		sendError(w, err.Error())
		return
	}

	// --- Business Logic ---
	rate, found := calculator.FindRate(req.State, req.Year, req.ProductCategory)
	if !found {
		msg := fmt.Sprintf("Tax rules for %s in %d are not available for the %s category.", req.State, req.Year, req.ProductCategory)
		sendError(w, msg)
		return
	}

	// --- Calculation ---
	result := calculator.Calculate(req.Amount, rate)

	// TODO: create a mapper function build TaxResponse with `req` and `result` values
	finalData := &TaxResponse{
		BaseAmount: req.Amount,
		TaxAmount:  result.TaxAmount,
		Total:      result.Total,
		Rate:       result.Rate,
		State:      req.State,
		Year:       req.Year,
	}

	// TODO: the sendError function also sends a JSON response, so create a helper function and re-use it in both places
	json.NewEncoder(w).Encode(APIResponse{
		Success: true,
		Data:    finalData,
	})
}

// --- Helpers and utilities ---
func setupCORS(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type")
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
