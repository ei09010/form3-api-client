# Go client for Form3 Organisation Account Services

### I don't have commercial experience with Go.

## Description

This is the Go Client for Form3 API Account Services. It is a Go library that enable a user to execute the following operations with the Form3 Accounts API:

- Create
- Fetch
- Delete

## Instructions

 - Install Docker's most recent version

 - Checkout this repo and run `docker-compose up` and watch the tests run!

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

