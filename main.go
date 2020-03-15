package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// Looking for a single argument at the command line: port #
	// If not specified, exit
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("Port not found. Syntax: \n\n\t quoteservice [port #]")
		os.Exit(1)
	}

	// Converting port # to int; exit if failure
	port, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatalf(fmt.Sprintf("Invalid port number: '%s'. Aborting.", args[0]))
		os.Exit(1)
	}

	app := &QuoteApplication{Quotes: QuoteList{}}

	app.Quotes.loadQuotes()

	http.HandleFunc("/", app.RootHandler)
	http.HandleFunc("/api/quote/", app.QuoteHandler)
	http.HandleFunc("/api/submit/", app.SubmitHandler)
	http.HandleFunc("/api/approve/", app.ApproveHandler)
	http.HandleFunc("/api/disapprove/", app.DisapproveHandler)

	log.Println(fmt.Sprintf("running on port %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
