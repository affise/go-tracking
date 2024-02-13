# go-affise-tracking [![GoDoc](https://godoc.org/github.com/affise/go-affise-tracking?status.svg)](https://godoc.org/github.com/affise/go-affise-tracking)

- [Clicks tracking](#clicks-tracking)
- [Conversions tracking](#conversions-tracking)


## Clicks tracking

To use first-party cookie there are two middlewares and two methods that's used by middlewares.

Middlewares
- **SetCookieMiddleware** is net/http middleware. It calls SetCookie to write click id cookie header and may throw ErrNoQueryParam error to channel.
- **MustSetCookieMiddleware** is net/http middleware. It calls MustSetCookie to write click id cookie header.

Methods
- **SetCookie** writes click id cookie to http.ResponseWriter and may throw ErrNoQueryParam error to channel.
- **MustSetCookie** writes click id cookie to http.ResponseWriter. This method ignores errors.

## Conversions tracking

To make s2s-integration you can use Postback settings and methods with it.

Postback has one required field - ClickID. Other fields are optional.

To make requests or do them it's useful create new PostbackProvider using NewPostbackProvider func.

This provider has next methods:
- **Request** returns http.Request for specified Postback.
- **RequestWithCookie** returns http.Request for specified Postback filling click id from specified http.Request.
- **Do** requests Postback with specified http.Client.
- **DoWithCookie** requests Postback with http.Client using click id cookie from http.Request.
- **DoDefault** requests Postback with http.DefaultClient.
- **DoDefaultWithCookie** requests Postback with http.DefaultClient using click id cookie from http.Request.
