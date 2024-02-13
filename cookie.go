package tracking

import (
	"errors"
	"fmt"
	"net/http"
	"time"
)

// These constants contain click id query param variations.
const (
	QueryParamVar1 = "click_id"
	QueryParamVar2 = "clickid"
	QueryParamVar3 = "afclick"
)

// CookieName is cookie name const.
const CookieName = "afclick"

// ErrNoQueryParam is error of request has no click id param.
var ErrNoQueryParam = errors.New("query has no click id param")

// SetCookieMiddleware is net/http middleware. It calls SetCookie to write click id cookie header.
func SetCookieMiddleware(next http.Handler, errs chan<- error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		SetCookie(w, r, errs)
		next.ServeHTTP(w, r)
	})
}

// MustSetCookieMiddleware is net/http middleware. It calls MustSetCookie to write click id cookie header.
func MustSetCookieMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		MustSetCookie(w, r)
		next.ServeHTTP(w, r)
	})
}

// SetCookie writes click id cookie to http.ResponseWriter.
// If there is no variation of query param it send error to chan.
func SetCookie(w http.ResponseWriter, r *http.Request, errs chan<- error) {
	var val string

	q := r.URL.Query()
	for _, n := range []string{QueryParamVar1, QueryParamVar2, QueryParamVar3} {
		if !q.Has(n) {
			continue
		}

		val = q.Get(n)

		break
	}

	if val == "" {
		if errs != nil {
			select {
			case <-r.Context().Done():
			case errs <- fmt.Errorf("%w: %s", ErrNoQueryParam, r.URL.String()):
			}
		}

		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:       CookieName,
		Value:      val,
		Path:       "/",
		Domain:     r.Host,
		Expires:    time.Now().Add(365 * 24 * time.Hour),
		RawExpires: "",
		MaxAge:     0,
		Secure:     true,
		HttpOnly:   false,
		SameSite:   http.SameSiteNoneMode,
		Raw:        "",
		Unparsed:   nil,
	})
}

// MustSetCookie writes click id cookie to http.ResponseWriter.
func MustSetCookie(w http.ResponseWriter, r *http.Request) {
	SetCookie(w, r, nil)
}
