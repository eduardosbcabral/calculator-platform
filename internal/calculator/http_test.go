package calculator

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		contentType string
		body        string
		wantStatus  int
		wantResult  float64
		wantCode    string
	}{
		{name: "binary operation", method: http.MethodPost, contentType: "application/json", body: `{"operation":"add","left":7,"right":5}`, wantStatus: http.StatusOK, wantResult: 12},
		{name: "zero operands are present", method: http.MethodPost, contentType: "application/json", body: `{"operation":"add","left":0,"right":0}`, wantStatus: http.StatusOK},
		{name: "unary operation", method: http.MethodPost, contentType: "application/json", body: `{"operation":"square_root","left":81}`, wantStatus: http.StatusOK, wantResult: 9},
		{name: "percentage", method: http.MethodPost, contentType: "application/json; charset=utf-8", body: `{"operation":"percentage","left":15,"right":200}`, wantStatus: http.StatusOK, wantResult: 30},
		{name: "missing left", method: http.MethodPost, contentType: "application/json", body: `{"operation":"add","right":5}`, wantStatus: http.StatusBadRequest, wantCode: "missing_operand"},
		{name: "missing right", method: http.MethodPost, contentType: "application/json", body: `{"operation":"add","left":5}`, wantStatus: http.StatusBadRequest, wantCode: "missing_operand"},
		{name: "unexpected right", method: http.MethodPost, contentType: "application/json", body: `{"operation":"square_root","left":81,"right":2}`, wantStatus: http.StatusBadRequest, wantCode: "unexpected_operand"},
		{name: "unsupported operation", method: http.MethodPost, contentType: "application/json", body: `{"operation":"modulo","left":7,"right":5}`, wantStatus: http.StatusBadRequest, wantCode: "unsupported_operation"},
		{name: "division by zero", method: http.MethodPost, contentType: "application/json", body: `{"operation":"divide","left":7,"right":0}`, wantStatus: http.StatusBadRequest, wantCode: "division_by_zero"},
		{name: "negative square root", method: http.MethodPost, contentType: "application/json", body: `{"operation":"square_root","left":-1}`, wantStatus: http.StatusBadRequest, wantCode: "negative_square_root"},
		{name: "unknown field", method: http.MethodPost, contentType: "application/json", body: `{"operation":"add","left":7,"right":5,"extra":true}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_json"},
		{name: "trailing object", method: http.MethodPost, contentType: "application/json", body: `{"operation":"add","left":7,"right":5}{}`, wantStatus: http.StatusBadRequest, wantCode: "invalid_json"},
		{name: "malformed json", method: http.MethodPost, contentType: "application/json", body: `{`, wantStatus: http.StatusBadRequest, wantCode: "invalid_json"},
		{name: "wrong content type", method: http.MethodPost, contentType: "text/plain", body: `{}`, wantStatus: http.StatusUnsupportedMediaType, wantCode: "unsupported_media_type"},
		{name: "wrong method", method: http.MethodGet, wantStatus: http.StatusMethodNotAllowed, wantCode: "method_not_allowed"},
		{name: "large request", method: http.MethodPost, contentType: "application/json", body: `{"operation":"` + strings.Repeat("x", maxRequestBytes) + `","left":7,"right":5}`, wantStatus: http.StatusRequestEntityTooLarge, wantCode: "request_too_large"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/api/calculate", strings.NewReader(test.body))
			if test.contentType != "" {
				request.Header.Set("Content-Type", test.contentType)
			}
			recorder := httptest.NewRecorder()

			handleCalculate(recorder, request)

			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", recorder.Code, test.wantStatus, recorder.Body.String())
			}
			if got := recorder.Header().Get("Content-Type"); got != "application/json" {
				t.Fatalf("Content-Type = %q, want application/json", got)
			}

			if test.wantStatus == http.StatusOK {
				var got response
				if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}
				if got.Result != test.wantResult {
					t.Fatalf("result = %v, want %v", got.Result, test.wantResult)
				}
				return
			}

			var got errorResponse
			if err := json.NewDecoder(recorder.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}
			if got.Error.Code != test.wantCode {
				t.Fatalf("error code = %q, want %q", got.Error.Code, test.wantCode)
			}
		})
	}
}

func TestAPIFallback(t *testing.T) {
	mux := http.NewServeMux()
	Register(mux)
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/unknown", nil))

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusNotFound)
	}
}
