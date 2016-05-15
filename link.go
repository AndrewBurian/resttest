package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Page struct {
	TotalCount   int           `json:"totalCount"`
	Page         int           `json:"page"`
	Transactions []Transaction `json:"transactions"`
}

type Transaction struct {
	Date     string `json:"Date"`
	DateTime time.Time
	Ledger   string  `json:"Ledger"`
	Amount   float32 `json:"Amount,string"`
	Company  string  `json:"Company"`
}

/*
First stage of the pipeline

Gets pages from API in sequence, does any intial processing,
and passes the resulting transactions down
*/
func Link(transactions chan<- Transaction) {

	transactionCount := 0

	for i := 1; ; i++ {

		// get page i json from API
		page, err := Get(i)

		// break on a non 200 response
		if err != nil {
			log.Println(err)
			break
		}

		// enrich the string data to a time.Time object
		for j, _ := range page.Transactions {
			DateTime(&page.Transactions[j])
		}

		// Send along all transactions
		for _, trans := range page.Transactions {
			transactions <- trans
			transactionCount++
		}

		if transactionCount == page.TotalCount {
			break
		}
	}

	close(transactions)
}

// Get a single page from the API
func Get(pageno int) (p Page, err error) {

	// Create the url for the page
	url := fmt.Sprintf("http://resttest.bench.co/transactions/%d.json", pageno)

	// Issue the GET request
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	// Check response code
	if resp.StatusCode != http.StatusOK {
		return p, fmt.Errorf("API responded with code %s", resp.Status)
	}

	// Load the data
	dat, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}

	// unmarshal data
	err = json.Unmarshal([]byte(dat), &p)
	if err != nil {
		return p, err
	}

	return p, nil
}

// Convert a datetime string into a time.Time object
func DateTime(t *Transaction) {

	// no error handing because POC
	// We assume the time data is well formed
	t.DateTime, _ = time.Parse("2006-01-02", t.Date)

}
