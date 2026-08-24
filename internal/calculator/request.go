package calculator

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
)

const maxRequestBytes = 4096

type calculationRequest struct {
	Operation operation `json:"operation"`
	Left      *float64  `json:"left"`
	Right     *float64  `json:"right"`
}

type requestError struct {
	status  int
	code    string
	message string
}

func decodeRequest(w http.ResponseWriter, r *http.Request) (calculationRequest, *requestError) {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		return calculationRequest{}, requestFailure(
			http.StatusUnsupportedMediaType,
			"unsupported_media_type",
			"Content-Type must be application/json.",
		)
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBytes)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	var input calculationRequest
	if err := decoder.Decode(&input); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			return calculationRequest{}, requestFailure(
				http.StatusRequestEntityTooLarge,
				"request_too_large",
				"Request body is too large.",
			)
		}
		return calculationRequest{}, invalidJSONError()
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return calculationRequest{}, invalidJSONError()
	}

	return input, nil
}

func invalidJSONError() *requestError {
	return requestFailure(
		http.StatusBadRequest,
		"invalid_json",
		"Request body must contain one valid JSON object with known fields.",
	)
}

func requestFailure(status int, code, message string) *requestError {
	return &requestError{status: status, code: code, message: message}
}
