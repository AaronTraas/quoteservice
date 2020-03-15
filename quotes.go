package main

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

func (quotes *QuoteList) setApproval(id int, approved bool) (*QuoteEntry, bool) {

	quoteIndex, quote := quotes.getQuoteById(id)
	if quoteIndex < 0 {
		return nil, false
	}

	quote.Approved = approved
	(*quotes)[quoteIndex] = *quote

	return quote, true
}

func (quotes *QuoteList) getQuoteById(id int) (int, *QuoteEntry) {
	for index, quote := range *quotes {
		if quote.Id == id {
			return index, &quote
		}
	}
	return -1, nil
}
