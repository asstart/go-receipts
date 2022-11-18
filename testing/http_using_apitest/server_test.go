package http_using_apitest_test

import (
	"testing"

	"github.com/asstart/go-receipts/testing/http"
	"github.com/steinfletcher/apitest"
)

func TestHandler(t *testing.T) {

}

func TestRouter(t *testing.T) {
	tt := []struct {
		name            string
		method          string
		queryParam      string
		expectedCode    int
		expectedMessage string
	}{
		{"GET /route?count=2", "GET", "2", 200, "Input: 2"},
		{"POST /route?count=2", "POST", "2", 405, "Method not allowed\n"},
		{"GET /route?count=", "GET", "", 400, "Query param not found\n"},
		{"GET /route?count=a", "GET", "a", 400, "Bad query param\n"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			apitest.New().
				Handler(http.Route()).
				Method(tc.method).
				URL("/route").
				Query("count", tc.queryParam).
				Expect(t).
				Status(tc.expectedCode).
				Body(tc.expectedMessage).
				End()
		})
	}
}
