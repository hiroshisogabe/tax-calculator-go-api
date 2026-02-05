package main

import (
	"testing"
)

func TestTaxRequest_Validate(t *testing.T) {
	tests := []struct {
		name          string
		input         TaxRequest
		expectedError bool
		expectedMsg   string
	}{
		{
			name: "Valid Input",
			input: TaxRequest{
				Amount:          100.0,
				State:           "ny",
				Year:            2024,
				ProductCategory: "electronics",
			},
			expectedError: false,
		},
		{
			name: "Invalid Amount (Zero)",
			input: TaxRequest{
				Amount:          0,
				State:           "NY",
				Year:            2024,
				ProductCategory: "electronics",
			},
			expectedError: true,
			expectedMsg:   "Amount must be greater than zero",
		},
		{
			name: "Invalid State (Too short)",
			input: TaxRequest{
				Amount:          100,
				State:           "N",
				Year:            2024,
				ProductCategory: "electronics",
			},
			expectedError: true,
			expectedMsg:   "State code is required (e.g., NY)",
		},
		{
			name: "Invalid Year (Out of range)",
			input: TaxRequest{
				Amount:          100,
				State:           "NY",
				Year:            202,
				ProductCategory: "electronics",
			},
			expectedError: true,
			expectedMsg:   "Year must be a 4-digit number",
		},
		{
			name: "Empty Category",
			input: TaxRequest{
				Amount:          100,
				State:           "NY",
				Year:            2024,
				ProductCategory: "",
			},
			expectedError: true,
			expectedMsg:   "Category is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We use a local copy to ensure normalization is harmless between tests
			req := tt.input
			err := req.Validate()

			if (err != nil) != tt.expectedError {
				t.Fatalf("Validate() error = %v, expectedError %v", err, tt.expectedError)
			}

			if tt.expectedError && err.Error() != tt.expectedMsg {
				t.Errorf("Validate() error message = %q, want %q", err.Error(), tt.expectedMsg)
			}

			// Check normalization for the happy path
			if !tt.expectedError && req.State != "NY" {
				t.Errorf("Validate() failed to normalize State: got %v, expected NY", req.State)
			}
		})
	}
}
