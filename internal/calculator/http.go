package calculator

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
)

const maxRequestBytes = 4096

type request struct {
	Operation operation `json:"operation"`
	Left      *float64  `json:"left"`
	Right     *float64  `json:"right"`
}

type response struct {
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error apiError `json:"error"`
}

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/calculate", handleCalculate)
	mux.HandleFunc("/api/", func(w http.ResponseWriter, _ *http.Request) {
		writeError(w, http.StatusNotFound, "not_found", "Endpoint not found.")
	})
}

func handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Use POST for this endpoint.")
		return
	}

	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(w, http.StatusUnsupportedMediaType, "unsupported_media_type", "Content-Type must be application/json.")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var input request
	if err := decoder.Decode(&input); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, "request_too_large", "Request body is too large.")
			return
		}
		writeError(w, http.StatusBadRequest, "invalid_json", "Request body must contain one valid JSON object with known fields.")
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		writeError(w, http.StatusBadRequest, "invalid_json", "Request body must contain one valid JSON object with known fields.")
		return
	}

	if input.Left == nil {
		writeError(w, http.StatusBadRequest, "missing_operand", "The left operand is required.")
		return
	}

	if input.Operation == squareRoot {
		if input.Right != nil {
			writeError(w, http.StatusBadRequest, "unexpected_operand", "Square root accepts only the left operand.")
			return
		}
	} else if isBinary(input.Operation) {
		if input.Right == nil {
			writeError(w, http.StatusBadRequest, "missing_operand", "The right operand is required.")
			return
		}
	} else {
		writeError(w, http.StatusBadRequest, "unsupported_operation", "Choose a supported operation.")
		return
	}

	right := 0.0
	if input.Right != nil {
		right = *input.Right
	}

	result, err := calculate(input.Operation, *input.Left, right)
	if err != nil {
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
		return
	}

	writeJSON(w, http.StatusOK, response{Result: result})
}

func isBinary(op operation) bool {
	switch op {
	case add, subtract, multiply, divide, power, percentage:
		return true
	default:
		return false
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
