package http_using_apitest

import (
	"fmt"
	"net/http"
	"strconv"
)

func run() {
	srv := Route()
	err := http.ListenAndServe(":8080", srv)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func Route() http.Handler {
	srv := http.NewServeMux()
	srv.HandleFunc("/route", Handler)
	return srv
}

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method

	if m != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	values := r.URL.Query()

	param := values.Get("count")
	if param == "" {
		http.Error(w, "Query param not found", http.StatusBadRequest)
		return
	}

	count, err := strconv.Atoi(param)
	if err != nil {
		http.Error(w, "Bad query param", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "Input: %v", count)
	w.WriteHeader(200)
}
