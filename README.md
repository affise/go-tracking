# go-tracking [![GoDoc](https://godoc.org/github.com/affise/go-tracking?status.svg)](https://godoc.org/github.com/affise/go-tracking)

Affise tracking SDK for go.

## Installation

To install it in the GOPATH:
```
go get https://github.com/affise/go-tracking
```
## Documentation

The links bellow should provide all the documentation needed to make the best
use of the library and the Segment API:

- [Documentation](https://help-center.affise.com/en/articles/6466563-postback-integration-s2s-admins)
- [godoc](https://godoc.org/github.com/affise/go-tracking)

## Usage

### Clicks

```go
package main

import (
	"log"
	"net/http"

	"github.com/affise/go-tracking"
)

func main() {
	http.HandleFunc("/click", func(w http.ResponseWriter, r *http.Request) {
		// request should contain param click_id/clickid/afclick 
		tracking.MustSetCookie(w, r)
		
		// ...
	})

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}
```

### Conversions

```go
package main

import (
	"log"
	"net/http"

	"github.com/affise/go-tracking"
)

func main() {
	pp := tracking.NewPostbackProvider("offers-client.affise.com", true)
	
	http.HandleFunc("/postback", func(w http.ResponseWriter, r *http.Request) {
		// request should contain first-party cookie 
		err := pp.DoDefaultWithCookie(r, &tracking.Postback{
			ActionID: "advertiser action id",
		})
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		
		// ...
	})

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatalf("server error: %v", err)
	}
}

```