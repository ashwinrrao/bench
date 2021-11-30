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
	pageOneResponse, err := getTransactions(1)
	if err != nil {
		fmt.Println("there was an error", err)
		return
	}

	// formula to calculate the total number of pages
	pages := int(math.Ceil(float64(pageOneResponse.TotalCount) / float64(len(pageOneResponse.Transactions))))

	// get the remaining pages asynchronously
	responses, err := getPagesAsync(pages)
	if err != nil {
		fmt.Println("there was an error", err)
		return
	}

	// append the transactions from the first page
	responses = append(responses, pageOneResponse)

	// calculate the running sum in O(n) time
	result := calcRunningSum(responses)

	// print out the result
	for key, val := range result {
		fmt.Println(key, val)
	}
}

func getTransactions(pageNumber int) (Response, error) {
	var response Response

	// perform the HTTP Get operation on the URL
	httpResponse, err := http.Get(fmt.Sprintf(URL, pageNumber))

	// if there was an error or if the status code is not 200, return an error
	if err != nil || httpResponse.StatusCode != http.StatusOK {
		return response, errGet
	}
	defer httpResponse.Body.Close()

	// decode the json body and catch any errors while decoding
	if err := json.NewDecoder(httpResponse.Body).Decode(&response); err != nil {
		return response, errMalformedResponse
	}

	return response, nil
}

func getPagesAsync(pages int) ([]Response, error) {
	var wg sync.WaitGroup
	responseChannel := make(chan AsyncResponse)
	var err error

	// loop through the second page to the last page
	// creating a thread for each page that Gets the data asynchronously
	for i := 2; i <= pages; i++ {
		wg.Add(1)
		go func(j int) {
			defer wg.Done()
			r, err := getTransactions(j)
			responseChannel <- AsyncResponse{r: r, err: err}
		}(i)
	}

	// close the channel once all the threads are finished
	go func() {
		wg.Wait()
		close(responseChannel)
	}()

	var result []Response
	// collect all the responses from all the threads and add them to a list
	for response := range responseChannel {
		result = append(result, response.r)
		if response.err != nil {
			err = response.err
		}
	}

	return result, err
}

func calcRunningSum(responses []Response) map[string]float64 {
	result := map[string]float64{}
	// loop through each transaction and update the running sum
	for _, response := range responses {
		for _, transaction := range response.Transactions {
			result[transaction.Date] += transaction.Amount
		}
	}
	return result
}
