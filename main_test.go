package main

import (
	"reflect"
	"testing"
)

func Test_getTransactions(t *testing.T) {
	type args struct {
		pageNumber int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "get transactions from the first page",
			args:    args{pageNumber: 1},
			want:    1, // hard coding the expected result
			wantErr: false,
		},
		{
			name:    "get transactions from the second page",
			args:    args{pageNumber: 2},
			want:    2, // hard coding the expected result
			wantErr: false,
		},
		{
			name:    "get transactions from the third page",
			args:    args{pageNumber: 3},
			want:    3, // hard coding the expected result
			wantErr: false,
		},
		{
			name:    "get transactions from the fourth page",
			args:    args{pageNumber: 4},
			want:    4, // hard coding the expected result
			wantErr: false,
		},
		{
			name:    "get transactions from a non-existent page",
			args:    args{pageNumber: 8},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTransactions(tt.args.pageNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Page, tt.want) {
				t.Errorf("getTransactions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPagesAsync(t *testing.T) {
	type args struct {
		pages int
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "when there is only page worth of transactions",
			args: args{
				pages: 1,
			},
			want:    0, // expecting to get no results, because this function excludes the first page
			wantErr: false,
		},
		{
			name: "when all pages are valid",
			args: args{
				pages: 4,
			},
			want:    3, // expecting to get 3 pages of results excluding the first page
			wantErr: false,
		},
		{
			name: "when some pages are invalid",
			args: args{
				pages: 8,
			},
			want:    7,    // expecting to get 7 pages of results excluding the first page
			wantErr: true, // also expecting an error because some pages gave a 404
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getPagesAsync(tt.args.pages)
			if (err != nil) != tt.wantErr {
				t.Errorf("getPagesAsync() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("getPagesAsync() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_calcRunningSum(t *testing.T) {
	type args struct {
		r []Response
	}
	tests := []struct {
		name string
		args args
		want map[string]float64
	}{
		{
			name: "when the input is empty",
			args: args{
				r: []Response{},
			},
			want: map[string]float64{}, // expect the output to be empty
		},
		{
			name: "when the input contains a single transaction",
			args: args{
				r: []Response{{Transactions: []Transaction{{Date: "2021-11-30", Amount: -11}}}},
			},
			want: map[string]float64{"2021-11-30": -11}, // expect the output to contain a single transaction
		},
		{
			name: "when the input contains a multiple transactions for a single day",
			args: args{
				r: []Response{{Transactions: []Transaction{{Date: "2021-11-30", Amount: -11}, {Date: "2021-11-30", Amount: -10}}}},
			},
			want: map[string]float64{"2021-11-30": -21}, // expect the output to contain a single entry with the daily sum
		},
		{
			name: "when the input contains a multiple transactions for multiple days",
			args: args{
				r: []Response{{Transactions: []Transaction{
					{Date: "2021-11-30", Amount: -11},
					{Date: "2021-11-30", Amount: -10},
					{Date: "2021-11-29", Amount: 15},
					{Date: "2021-11-29", Amount: 25},
					{Date: "2021-11-28", Amount: -14.5}}}},
			},
			want: map[string]float64{"2021-11-30": -21, "2021-11-29": 40, "2021-11-28": -14.5}, // expect the output to contain a multiple entries with the daily sum
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcRunningSum(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calcRunningSum() = %v, want %v", got, tt.want)
			}
		})
	}
}
