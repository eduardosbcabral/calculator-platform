package calculator

import "net/http"

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

	input, requestError := decodeRequest(w, r)
	if requestError != nil {
		writeError(w, requestError.status, requestError.code, requestError.message)
		return
	}

	left, right, requestError := input.operands()
	if requestError != nil {
		writeError(w, requestError.status, requestError.code, requestError.message)
		return
	}

	result, err := calculate(input.Operation, left, right)
	if err != nil {
		writeCalculationError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, calculationResponse{Result: result})
}
