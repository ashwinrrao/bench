package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sync"
)

const (
	URL = "https://resttest.bench.co/transactions/%v.json"
)

var (
	errGet               = errors.New("there was an error getting the data")
	errMalformedResponse = errors.New("the data returned is in an invalid format")
)

func main() {

	// get the first page of transactions
	r, err := getTransactions(1)
	if err != nil {
		fmt.Println("there was an error", err)
		return
	}

	// formula to calculate the total number of pages
	pages := int(math.Ceil(float64(r.TotalCount) / float64(len(r.Transactions))))

	// get the remaining pages asynchronously
	resp, err := getPagesAsync(pages)
	if err != nil {
		fmt.Println("there was an error", err)
		return
	}

	// append the transactions from the first page
	resp = append(resp, r)

	// calculate the running sum in O(n) time
	result := calcRunningSum(resp)

	// print out the result
	for key, val := range result {
		fmt.Println(key, val)
	}
}

func getTransactions(pageNumber int) (Response, error) {
	var r Response

	// perform the HTTP Get operation on the URL
	resp, err := http.Get(fmt.Sprintf(URL, pageNumber))

	// if there was an error or if the status code is not 200, return an error
	if err != nil || resp.StatusCode != http.StatusOK {
		return r, errGet
	}
	defer resp.Body.Close()

	// decode the json body and catch any errors while decoding
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return r, errMalformedResponse
	}

	return r, nil
}

func getPagesAsync(pages int) ([]Response, error) {
	var wg sync.WaitGroup
	respChan := make(chan AsyncResponse)
	var err error

	// loop through the second page to the last page
	// creating a thread for each page that Gets the data asynchronously
	for i := 2; i <= pages; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			r, err := getTransactions(j)
			respChan <- AsyncResponse{r: r, err: err}
		}(i)
	}

	// close the channel once all the threads are finished
	go func() {
		wg.Wait()
		close(respChan)
	}()

	var r []Response
	// collect all the responses from all the threads and add them to a list
	for resp := range respChan {
		r = append(r, resp.r)
		if resp.err != nil {
			err = resp.err
		}
	}

	return r, err
}

func calcRunningSum(r []Response) map[string]float64 {
	result := map[string]float64{}
	// loop through each transaction and update the running sum
	for _, resp := range r {
		for _, t := range resp.Transactions {
			result[t.Date] += t.Amount
		}
	}
	return result
}
