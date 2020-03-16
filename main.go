package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"strconv"
)

func (app *QuoteApplication) dumpQuotes() {
	if app.QuoteDbJsonPath == "" {
		return
	}

	writer, err := os.OpenFile(app.QuoteDbJsonPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Println(err)
		return
	}

	dumpQuotesToWriter(app, writer)
	writer.Close()
}

func loadQuotes(app *QuoteApplication) {

	reader, err := os.Open(app.QuoteDbJsonPath)
	if err != nil {
		log.Println(err)
		return
	}

	loadQuotesFromReader(app, reader)
	reader.Close()
}

func main() {
	// Looking for a single argument at the command line: port #
	// If not specified, exit
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatalf("Port not found. Syntax: \n\n\t quoteservice [port #]")
	}

	// Converting port # to int; exit if failure
	port, err := strconv.ParseUint(args[0], 10, 16)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Invalid port number: '%s'. Aborting.", args[0]))
	}

	usr, _ := user.Current()
	app := &QuoteApplication{
		Quotes:          QuoteList{},
		QuoteDbJsonPath: usr.HomeDir + "/.quoteservice.json",
	}

	loadQuotes(app)

	http.HandleFunc("/", app.RootHandler)
	http.HandleFunc("/api/quote/", app.QuoteHandler)
	http.HandleFunc("/api/submit/", app.SubmitHandler)
	http.HandleFunc("/api/approve/", app.ApproveHandler)
	http.HandleFunc("/api/disapprove/", app.DisapproveHandler)

	log.Println(fmt.Sprintf("running on port %d", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
