package calculator

import (
	"encoding/json"
	"errors"
	"net/http"
)

type calculationResponse struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeCalculationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errDivisionByZero):
		writeError(w, http.StatusBadRequest, "division_by_zero", "Cannot divide by zero.")
	case errors.Is(err, errNegativeSquareRoot):
		writeError(w, http.StatusBadRequest, "negative_square_root", "Cannot calculate the square root of a negative number.")
	case errors.Is(err, errNonFiniteResult):
		writeError(w, http.StatusBadRequest, "non_finite_result", "The result is outside the supported numeric range.")
	default:
		writeError(w, http.StatusBadRequest, "unsupported_operation", "Choose a supported operation.")
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: apiError{Code: code, Message: message}})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
