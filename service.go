package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Response struct {
	Path    string      `json:"path"`
	Success bool        `json:"success"`
	Status  int         `json:"status"`
	Quote   *QuoteEntry `json:"quote,omitempty"`
	Quotes  *QuoteList  `json:"quotes,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type QuoteApplication struct {
	Quotes QuoteList
}

func sendJsonResponse(w http.ResponseWriter, response Response) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if response.Status == 0 {
		response.Status = 200
	}

	w.WriteHeader(response.Status)

	responseJson, _ := json.MarshalIndent(response, "", "  ")

	fmt.Fprintf(w, string(responseJson))
}

func setApprovalResponse(quotes *QuoteList, path string, idString string, approved bool) Response {
	if idString == "" {
		return Response{
			Path:    path,
			Success: false,
			Status:  http.StatusBadRequest,
			Error:   "'id' is required parameter",
		}
	}
	id, errIdConvert := strconv.Atoi(idString)
	if errIdConvert != nil {
		return Response{
			Path:    path,
			Success: false,
			Status:  http.StatusBadRequest,
			Error:   "'id' is not an integer",
		}
	}

	quote, success := quotes.setApproval(id, approved)
	if success == false {
		return Response{
			Path:    path,
			Success: false,
			Status:  http.StatusNotFound,
			Error:   fmt.Sprintf("quote #%d not found", id),
		}
	}

	return Response{
		Path:    path,
		Success: true,
		Quote:   quote,
	}
}

func (app *QuoteApplication) rootHandler(w http.ResponseWriter, r *http.Request) {
	var response Response

	if r.Method == http.MethodGet {
		response = Response{
			Path:    r.URL.Path,
			Success: true,
			Quotes:  &app.Quotes,
		}
	}

	sendJsonResponse(w, response)
}

func (app *QuoteApplication) submitHandler(w http.ResponseWriter, r *http.Request) {
	bodyString := r.URL.Query().Get("body")
	urlString := r.URL.Query().Get("url")

	if (bodyString == "") || (urlString == "") {
		sendJsonResponse(w, Response{
			Path:    r.URL.Path,
			Status:  http.StatusBadRequest,
			Success: false,
			Error:   "'body' and 'url' are require parameters",
		})
		return
	}

	quote := QuoteEntry{
		Id:       app.Quotes.getNextId(),
		Body:     bodyString,
		Url:      urlString,
		Approved: false,
	}

	app.Quotes = append(app.Quotes, quote)

	sendJsonResponse(w, Response{
		Path:    r.URL.Path,
		Status:  http.StatusCreated,
		Success: true,
		Quote:   &quote,
	})
}

func (app *QuoteApplication) approveHandler(w http.ResponseWriter, r *http.Request) {
	sendJsonResponse(w, setApprovalResponse(&app.Quotes, r.URL.Path, r.URL.Query().Get("id"), true))
}

func (app *QuoteApplication) disapproveHandler(w http.ResponseWriter, r *http.Request) {
	sendJsonResponse(w, setApprovalResponse(&app.Quotes, r.URL.Path, r.URL.Query().Get("id"), false))
}
