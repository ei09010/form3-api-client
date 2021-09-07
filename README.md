# Go client for Form3 Organisation Account Services

### I don't have commercial experience with Go.

## Description

This is the Go Client for Form3 API Account Services. It is a Go library that enable a user to execute the following operations with the Form3 Accounts API:

- Create
- Fetch
- Delete

## Instructions

 - Install Docker's most recent version

 - Checkout this repo, run `docker-compose up` and watch the tests run!

## Usage

Sample usage of the Accounts API with a 10 second timeout value:

```go
package main

import (
	"ei09010/form3-api-client/organisation/accounts"
)

func main() {

	accountsClient, err := accounts.NewClient(accounts.WithTimeout(time.Duration(10 * time.Second)))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	accountResponse, err := accountsClient.Fetch(context.Background(), uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	fmt.Println(accountResponse)
}
```

## Production client nice to haves

- Rate limiting would maintain a predictable and safe maximum threshold on the amount of requests executed per second (or other given time interval)

- Connection re-usage between http requests for efficient resource usage ( both client and server side)

- Validators for the account object properties (in the Create method) could save unnecessary requests

- Monitoring: expose methods to make metrics available to be exported (Examples: requests per second, latency, error rate)