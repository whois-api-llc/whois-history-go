package whoishistory

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const defaultHistoricWhoisURL = `https://whois-history.whoisxmlapi.com/api/v1`

// HistoricService is an interface for Historic Whois API
type HistoricService interface {
	Purchase(ctx context.Context, name string, opts ...Option) ([]*WhoisRecord, *Response, error)
	Preview(ctx context.Context, name string, opts ...Option) (int, *Response, error)
}

type historicServiceOp struct {
	client  *Client
	baseURL *url.URL
}

var _ HistoricService = &historicServiceOp{}

func (service *historicServiceOp) newRequest() (*http.Request, error) {

	u, _ := url.Parse(service.baseURL.String())

	req, err := service.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("outputFormat", "JSON")
	query.Set("apiKey", service.client.apiKey)

	req.URL.RawQuery = query.Encode()

	return req, nil
}

type historicResponse struct {
	RecordsCount int            `json:"recordsCount"`
	Records      []*WhoisRecord `json:"records"`
	Code         int            `json:"code"`
	Message      string         `json:"messages"`
}

func (service *historicServiceOp) request(ctx context.Context, purchase bool, name string, opts ...Option) (*historicResponse, *Response, error) {
	if name == "" {
		return nil, nil, &ArgError{"name", "cannot be empty"}
	}

	req, err := service.newRequest()
	if err != nil {
		return nil, nil, err
	}

	q := req.URL.Query()
	q.Set("domainName", name)
	if purchase {
		q.Set("mode", "purchase")
	} else {
		q.Set("mode", "preview")
	}

	for _, opt := range opts {
		opt(q)
	}

	req.URL.RawQuery = q.Encode()

	var b strings.Builder
	resp, err := service.client.Do(ctx, req, &b)
	if err != nil {
		return nil, resp, err
	}

	respErr := checkResponse(resp.Response)

	response := historicResponse{}

	err = json.NewDecoder(strings.NewReader(b.String())).Decode(&response)
	if err != nil {
		if respErr != nil {
			return nil, resp, err
		}
		return nil, resp, fmt.Errorf("cannot parse response: %w", err)
	}

	if response.Message != "" || response.Code != 0 {
		return nil, resp, ErrorMessage{
			Code:    response.Code,
			Message: response.Message,
		}
	}

	if respErr != nil {
		return nil, resp, err
	}

	return &response, resp, nil
}

// Purchase returns the slice of records.
func (service *historicServiceOp) Purchase(ctx context.Context, name string, opts ...Option) ([]*WhoisRecord, *Response, error) {

	response, resp, err := service.request(ctx, true, name, opts...)
	if err != nil {
		return nil, resp, err
	}

	return response.Records, resp, nil
}

// Preview returns the number of records. No credits deducted.
func (service *historicServiceOp) Preview(ctx context.Context, name string, opts ...Option) (int, *Response, error) {

	response, resp, err := service.request(ctx, false, name, opts...)
	if err != nil {
		return 0, resp, err
	}

	return response.RecordsCount, resp, nil
}

// ArgError is an argument error
type ArgError struct {
	Name    string
	Message string
}

func (a *ArgError) Error() string {
	return `invalid argument: "` + a.Name + `" ` + a.Message
}
