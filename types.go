package main

type Transaction struct {
	Date    string  `json:"Date"`
	Ledger  string  `json:"Ledger"`
	Amount  float64 `json:"Amount,string"`
	Company string  `json:"Company"`
}

type Response struct {
	TotalCount   int           `json:"totalCount"`
	Page         int           `json:"page"`
	Transactions []Transaction `json:"transactions"`
}

type AsyncResponse struct {
	r   Response
	err error
}
