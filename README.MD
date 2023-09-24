# SPARING Go API

## Description

This is a simple SPARING api client. It is written in Go.

## Installation

```bash
 go get github.com/gamanlab/sparing
```

## Usage

Get secret key from SPARING server

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gamanlab/sparing"
)

func main() {
    api = NewSparingAPI("", svr.URL+"/secret", "")

    secret, err := api.GetSecret()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(secret)
}

```


Submit data to SPARING server

```go

package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gamanlab/sparing"
)

func main() {
    ph := float32(18692)
	cod := float32(5508)
	tss := float32(5466)
	nh3n := float32(16539)
	debit := float32(17006)

	payload := ApiPayload{
		UID:      1120800300014,
		DateTime: 1568630149,
		PH:       &ph,
		COD:      &cod,
		TSS:      &tss,
		NH3N:     &nh3n,
		Debit:    &debit,
	}


    api = NewSparingAPI("", svr.URL+"/submit", "")

    err := api.Submit(secret, payload)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("OK")
}

```


# License

MIT License

