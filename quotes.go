package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
)

type QuoteEntry struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	Url      string `json:"url"`
	Approved bool   `json:"approved"`
}

type QuoteList []QuoteEntry

func (quotes *QuoteList) getNextId() int {
	maxId := 0
	for _, quote := range *quotes {
		if quote.Id > maxId {
			maxId = quote.Id
		}
	}
	return maxId + 1
}

func (quotes *QuoteList) setApproval(path string, idString string, approved bool) Response {
	if idString == "" {
		return Response{
			Path:    path,
			Success: false,
			Error:   "'id' is required parameter",
		}
	}
	id, errIdConvert := strconv.Atoi(idString)
	if errIdConvert != nil {
		return Response{
			Path:    path,
			Success: false,
			Error:   "'id' is not an integer",
		}
	}

	quoteIndex, quote := quotes.getQuoteById(id)
	if quoteIndex < 0 {
		return Response{
			Path:    path,
			Success: false,
			Error:   fmt.Sprintf("quote #%d not found", id),
		}
	}

	quote.Approved = approved
	(*quotes)[quoteIndex] = *quote

	errDump := quotes.dumpQuotes()
	errDumpString := ""
	if errDump != nil {
		errDumpString = errDump.Error()
	}

	return Response{
		Path:    path,
		Quote:   quote,
		Success: errDump != nil,
		Error:   errDumpString,
	}
}

func (quotes *QuoteList) getQuoteById(id int) (int, *QuoteEntry) {
	for index, quote := range *quotes {
		if quote.Id == id {
			return index, &quote
		}
	}
	return -1, nil
}

func getQuoteDbPath() string {
	usr, _ := user.Current()
	fileName := usr.HomeDir + "/.quoteservice.json"

	return fileName
}

func (quotes *QuoteList) dumpQuotes() error {
	quoteListJson, _ := json.MarshalIndent(quotes, "", "  ")

	fileName := getQuoteDbPath()
	log.Println("Writing output to: ", fileName)

	_, errCreate := os.Create(fileName)
	if errCreate != nil {
		return errCreate
	}

	errWrite := ioutil.WriteFile(fileName, quoteListJson, 0644)
	return errWrite
}

func (quotes *QuoteList) loadQuotes() {
	fileName := getQuoteDbPath()

	log.Println("Loading DB from: ", fileName)

	if _, errExists := os.Stat("fileName"); errExists == nil {
		log.Println(errExists)
		return
	}

	quoteBytes, errRead := ioutil.ReadFile(fileName) // just pass the file name
	if errRead != nil {
		log.Println(errRead)
		return
	}

	errJson := json.Unmarshal(quoteBytes, quotes)
	if errJson != nil {
		log.Println(errJson)
	}
}
