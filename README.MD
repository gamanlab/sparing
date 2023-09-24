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

```
The MIT License (MIT)
=====================
SPARING Go API

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.