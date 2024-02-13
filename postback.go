package tracking

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// PostbackStatus is alias for uint8 to determine integer status value from other int values.
type PostbackStatus uint8

// PostbackStatus enum.
const (
	PostbackStatusInvalid PostbackStatus = iota
	PostbackStatusConfirmed
	PostbackStatusPending
	PostbackStatusDeclined
	_
	PostbackStatusHold
)

// CustomFieldsCount is amount of custom fields.
// Custom field 8-15 are allowed only if customer has accordant plan.
const CustomFieldsCount = 15

// URL paths.
const (
	pathPostback = "postback"
)

// URL params.
const (
	clickID      = "click_id"
	actionID     = "action_id"
	goal         = "goal"
	sum          = "sum"
	ip           = "ip"
	status       = "status"
	referrer     = "referrer"
	comment      = "comment"
	secure       = "secure"
	fbclid       = "fbclid"
	deviceType   = "device_type"
	userID       = "user_id"
	customFields = "custom_field%d"
)

var (
	// ErrInvalidClickID occurs when postback has invalid click id param.
	ErrInvalidClickID = errors.New("invalid click id")
	// ErrInvalidResponseStatus occurs when postback returns invalid response code.
	ErrInvalidResponseStatus = errors.New("invalid response status")
)

// Postback is setting of server postback. Only ClickID is required.
type Postback struct {
	ClickID      string
	ActionID     string
	Goal         string
	Sum          float64
	IP           *net.IP
	Status       PostbackStatus
	Referrer     string
	Comment      string
	Secure       string
	FbClID       string
	DeviceType   string
	UserID       string
	CustomFields [CustomFieldsCount]string
}

// PostbackProvider realizes postback sending functions.
type PostbackProvider struct {
	u *url.URL
}

// NewPostbackProvider return new PostbackProvider.
func NewPostbackProvider(domain string, ssl bool) *PostbackProvider {
	u := &url.URL{Scheme: "http", Host: domain, Path: pathPostback} //nolint:exhaustivestruct,exhaustruct
	if ssl {
		u.Scheme = "https"
	}

	return &PostbackProvider{u: u}
}

// Request returns http.Request for specified Postback.
// It may occur ErrInvalidClickID, syntax.Error or http-package error.
func (p *PostbackProvider) Request(ctx context.Context, pb *Postback) (*http.Request, error) {
	ok, err := regexp.Match("[0-9a-f]{24}", []byte(pb.ClickID))
	if err != nil {
		return nil, fmt.Errorf("failed to validate click id: %w", err)
	}

	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrInvalidClickID, pb.ClickID)
	}

	r, err := http.NewRequestWithContext(ctx, http.MethodGet, p.u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make new request: %w", err)
	}

	r.URL.RawQuery = p.query(pb).Encode()

	return r, nil
}

func (p *PostbackProvider) query(pb *Postback) url.Values {
	q := make(url.Values)
	q.Set(clickID, pb.ClickID)

	for k, v := range map[string]string{
		actionID:   pb.ActionID,
		goal:       pb.Goal,
		referrer:   pb.Referrer,
		comment:    pb.Comment,
		secure:     pb.Secure,
		fbclid:     pb.FbClID,
		deviceType: pb.DeviceType,
		userID:     pb.UserID,
	} {
		if v != "" {
			q.Set(k, v)
		}
	}

	if pb.Sum != 0 {
		q.Set(sum, fmt.Sprintf("%v", pb.Sum))
	}

	if pb.IP != nil {
		q.Set(ip, pb.IP.String())
	}

	//nolint:exhaustive
	switch pb.Status {
	case PostbackStatusConfirmed, PostbackStatusPending, PostbackStatusDeclined, PostbackStatusHold:
		q.Set(status, strconv.Itoa(int(pb.Status)))
	default:
	}

	for i := range pb.CustomFields {
		if pb.CustomFields[i] != "" {
			q.Set(fmt.Sprintf(customFields, i+1), pb.CustomFields[i])
		}
	}

	return q
}

// RequestWithCookie returns http.Request for specified Postback. It fills click id from specified http.Request.
// Addition to cases of Request error RequestWithCookie may occur ErrInvalidClickID when no click id cookie is set.
func (p *PostbackProvider) RequestWithCookie(r *http.Request, pb *Postback) (*http.Request, error) {
	c, err := r.Cookie(CookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, ErrInvalidClickID
		}

		return nil, fmt.Errorf("failed to get click id cookie: %w", err)
	}

	pb.ClickID = c.Value

	return p.Request(r.Context(), pb)
}

// Do requests Postback with specified http.Client.
func (p *PostbackProvider) Do(ctx context.Context, client *http.Client, pb *Postback) error {
	req, err := p.Request(ctx, pb)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}

	return p.do(client, req)
}

// DoWithCookie requests Postback with http.Client using click id cookie from http.Request.
func (p *PostbackProvider) DoWithCookie(r *http.Request, client *http.Client, pb *Postback) error {
	req, err := p.RequestWithCookie(r, pb)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}

	return p.do(client, req)
}

// DoDefault requests Postback with http.DefaultClient.
func (p *PostbackProvider) DoDefault(ctx context.Context, pb *Postback) error {
	return p.Do(ctx, http.DefaultClient, pb)
}

// DoDefaultWithCookie requests Postback with http.DefaultClient using click id cookie from http.Request.
func (p *PostbackProvider) DoDefaultWithCookie(r *http.Request, pb *Postback) error {
	return p.DoWithCookie(r, http.DefaultClient, pb)
}

func (p *PostbackProvider) do(client *http.Client, r *http.Request) error {
	res, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("failed to do request %s: %w", r.URL.String(), err)
	}

	if res.StatusCode != http.StatusOK {
		return ErrInvalidResponseStatus
	}

	if err := res.Body.Close(); err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}

	return nil
}
