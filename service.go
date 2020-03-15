package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
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
		response.Status = http.StatusOK
	}

	w.WriteHeader(response.Status)

	responseJson, _ := json.MarshalIndent(response, "", "  ")

	fmt.Fprintf(w, string(responseJson))
}

func setApprovalResponse(quotes *QuoteList, r *http.Request, approved bool) Response {
	if r.Method != http.MethodPost {
		return Response{
			Success: false,
			Status:  http.StatusMethodNotAllowed,
			Error:   fmt.Sprintf("Method %s not allowed. Must be POST.", r.Method),
		}
	}

	r.ParseForm()
	idString := r.Form.Get("id")

	if idString == "" {
		return Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Error:   "'id' is required parameter",
		}
	}
	id, errIdConvert := strconv.Atoi(idString)
	if errIdConvert != nil {
		return Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Error:   "'id' is not an integer",
		}
	}

	quote, success := quotes.setApproval(id, approved)
	if success == false {
		return Response{
			Success: false,
			Status:  http.StatusNotFound,
			Error:   fmt.Sprintf("quote #%d not found", id),
		}
	}

	return Response{
		Success: true,
		Quote:   quote,
	}
}

func (app *QuoteApplication) RootHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		sendJsonResponse(w, Response{
			Success: false,
			Status:  http.StatusMethodNotAllowed,
			Error:   fmt.Sprintf("Method %s not allowed. Must be GET.", r.Method),
		})
		return
	}

	sendJsonResponse(w, Response{
		Success: true,
		Quotes:  &app.Quotes,
	})
}

func (app *QuoteApplication) QuoteHandler(w http.ResponseWriter, r *http.Request) {

	if (r.Method != http.MethodGet) && (r.Method != http.MethodDelete) {
		sendJsonResponse(w, Response{
			Success: false,
			Status:  http.StatusMethodNotAllowed,
			Error:   fmt.Sprintf("Method %s not allowed. Must be GET or DELETE.", r.Method),
		})
		return
	}

	uriSegments := strings.Split(r.URL.Path, "/")

	idString := uriSegments[3]
	id, errIdConvert := strconv.Atoi(idString)
	if errIdConvert != nil {
		sendJsonResponse(w, Response{
			Success: false,
			Status:  http.StatusBadRequest,
			Error:   fmt.Sprintf("id '%s' is not an integer", idString),
		})
		return
	}

	quoteIndex, quote := app.Quotes.getQuoteById(id)
	if quoteIndex < 0 {
		sendJsonResponse(w, Response{
			Success: false,
			Status:  http.StatusNotFound,
			Error:   fmt.Sprintf("quote #%d not found", id),
		})
		return
	}

	sendJsonResponse(w, Response{
		Success: true,
		Quote:   quote,
	})
}

func (app *QuoteApplication) SubmitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJsonResponse(w, Response{
			Success: false,
			Status:  http.StatusMethodNotAllowed,
			Error:   fmt.Sprintf("Method %s not allowed. Must be POST.", r.Method),
		})
		return
	}

	r.ParseForm()
	bodyString := r.Form.Get("body")
	urlString := r.Form.Get("url")

	if (bodyString == "") || (urlString == "") {
		sendJsonResponse(w, Response{
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
		Status:  http.StatusCreated,
		Success: true,
		Quote:   &quote,
	})
}

func (app *QuoteApplication) ApproveHandler(w http.ResponseWriter, r *http.Request) {
	sendJsonResponse(w, setApprovalResponse(&app.Quotes, r, true))
}

func (app *QuoteApplication) DisapproveHandler(w http.ResponseWriter, r *http.Request) {
	sendJsonResponse(w, setApprovalResponse(&app.Quotes, r, false))
}
