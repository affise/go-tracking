# go-affise-tracking [![GoDoc](https://godoc.org/github.com/affise/go-affise-tracking?status.svg)](https://godoc.org/github.com/affise/go-affise-tracking)

1. [Clicks tracking](#clicks-tracking)


## Clicks tracking

To use first-party cookie there are two middlewares and two methods that's used by middlewares.

Middlewares
- **SetCookieMiddleware** is net/http middleware. It calls SetCookie to write click id cookie header and may throw ErrNoQueryParam error to channel.
- **MustSetCookieMiddleware** is net/http middleware. It calls MustSetCookie to write click id cookie header.

Methods
- **SetCookie** writes click id cookie to http.ResponseWriter and may throw ErrNoQueryParam error to channel.
- **MustSetCookie** writes click id cookie to http.ResponseWriter. This method ignores errors.

