# Go client for Form3 Organisation Account Services

### I don't have commercial experience with Go.

## Description

This the Go Client for Form3 API Account Services. It is a Go librar that enable a user to execute the following operations with the Form3 Accounts API:

- Create
- Fetch
- Delete

## Requirements

 - If you only whish to run the automated tests (unit, mock and integration) : Docker recent version
 - For development usage: Go 1.16.6 or later.

## Instructions

 - If you are a reviewer from form3tech-interviewer-1, just checkout this repo and run `docker-compose up`

## Instalation

To install the Go Client for Form3 Account Services, please execute the following `go get` command.

```bash
    go get ei09010/form3-api-client/organisation/accounts
```

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

	accountResponse, _, err := accountsClient.Fetch(context.Background(), uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"))
	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	fmt.Println(accountResponse)
}
```

