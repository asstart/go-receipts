package http_test

import (
	"fmt"
	gohttp "net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asstart/go-receipts/testing/http"
)

func TestHandler(t *testing.T) {

	tt := []struct {
		name            string
		method          string
		queryParam      string
		expectedCode    int
		expectedMessage string
	}{
		{"GET /route?count=2", "GET", "2", 200, "Input: 2"},
		{"POST /route?count=2", "POST", "2", 405, "Method not allowed"},
		{"GET /route?count=", "GET", "", 400, "Query param not found"},
		{"GET /route?count=a", "GET", "a", 400, "Bad query param"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest(tc.method, fmt.Sprintf("localhost:8080/route?count=%v", tc.queryParam), nil)
			w := httptest.NewRecorder()

			http.Handler(w, r)

		
			if w.Code != tc.expectedCode {
				t.Fatalf("expected: %v, actual: %v", tc.expectedCode, w.Code)
			}

			body := w.Body.String()
			if strings.TrimSpace(body) != tc.expectedMessage {
				t.Fatalf("expected body: %v, actual: %v", tc.expectedMessage, body)
			}
		})
	}
}

func TestRouter(t *testing.T) {
	srv := httptest.NewServer(http.Route())
	defer srv.Close()

	res, err := gohttp.Get(fmt.Sprintf("%s/route?count=2", srv.URL))
	if err != nil {
		t.Fatalf("expected success request, but got: %v", err)
	}
	if res.StatusCode != gohttp.StatusOK {
		t.Fatalf("expected %v:, actual: %v", gohttp.StatusOK, res.StatusCode)
	}
}