package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Path    string      `json:"path"`
	Success bool        `json:"success"`
	Quote   *QuoteEntry `json:"quote,omitempty"`
	Quotes  *QuoteList  `json:"quotes,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type QuoteApplication struct {
	Quotes QuoteList
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST")
	(*w).Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-Requested-With, Authorization, X-CSRF-Token, access-control-allow-origin, access-control-allow-headers")
	(*w).Header().Set("Content-Type", "application/json; charset=utf-8")
}

func (app *QuoteApplication) rootHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	response := Response{
		Path:    r.URL.Path,
		Success: true,
		Quotes:  &app.Quotes,
	}

	responseJson, _ := json.MarshalIndent(response, "", "  ")

	fmt.Fprintf(w, string(responseJson))
}

func (app *QuoteApplication) submitHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	bodyString := r.URL.Query().Get("body")
	urlString := r.URL.Query().Get("url")

	if (bodyString == "") || (urlString == "") {
		responseJson, _ := json.MarshalIndent(Response{
			Path:    r.URL.Path,
			Success: false,
			Error:   "'body' and 'url' are require parameters",
		}, "", "  ")

		fmt.Fprintf(w, string(responseJson))
		return
	}

	quote := QuoteEntry{
		Id:       app.Quotes.getNextId(),
		Body:     bodyString,
		Url:      urlString,
		Approved: false,
	}

	app.Quotes = append(app.Quotes, quote)

	responseJson, _ := json.MarshalIndent(Response{
		Path:    r.URL.Path,
		Success: true,
		Quote:   &quote,
	}, "", "  ")

	fmt.Fprintf(w, string(responseJson))
}

func (app *QuoteApplication) approveHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	idString := r.URL.Query().Get("id")

	response := app.Quotes.setApproval(r.URL.Path, idString, true)

	responseJson, _ := json.MarshalIndent(response, "", "  ")

	fmt.Fprintf(w, string(responseJson))
}

func (app *QuoteApplication) disapproveHandler(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)

	idString := r.URL.Query().Get("id")

	response := app.Quotes.setApproval(r.URL.Path, idString, false)

	responseJson, _ := json.MarshalIndent(response, "", "  ")

	fmt.Fprintf(w, string(responseJson))
}
