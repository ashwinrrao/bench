
# Bench Rest Test

This is a command line application that gets transactions from a URL and calculates the daily running sum for each day and displays it to the console.



## Running the app
**Note**: Open the command line and navigate to the folder containing the code before running the app or the tests

go run main.go types.go
## Running tests
go test

## Trade Offs
- The getPagesAsync method can be improved to get the first page of the transactions too. This would increase code modularity and readability. I decided to not fetch the first page again purely for efficiency reasons.
- Because I am using a hash map to store the running sum for each day, when I print the results the entries will be printed in random order instead of on a sequential basis.
- Because I am using GoLang there is a rounding issue when I add two floating point numbers. If I was doing this in production code, I would probably use a helper like the [Decimal library](https://pkg.go.dev/github.com/shopspring/decimal#section-readme)
## Authors

- [@ashwinrrao](https://www.github.com/ashwinrrao)