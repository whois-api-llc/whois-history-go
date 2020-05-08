package whoishistory

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const (
	libraryVersion = "0.1.0"
	userAgent      = "whoisxmlapi-go/" + libraryVersion
	mediaType      = "application/json"
)

// ClientParams is used to create Client. None of parameters are mandatory and
// leaving this struct empty works just fine for most cases.
type ClientParams struct {
	// HTTPClient is the client used to access API endpoint.
	// If it's nil then value API client uses http.DefaultClient
	HTTPClient *http.Client
	// Endpoint for `whois` service.
	WhoisBaseURL *url.URL
	// Endpoint for `historic whois` service.
	HistoricBaseURL *url.URL
}

// NewBasicClient creates Client with recommended parameters.
func NewBasicClient(APIKey string) *Client {
	return NewClient(APIKey, ClientParams{})
}

// NewClient creates Client with provided parameters.
func NewClient(apiKey string, params ClientParams) *Client {

	var err error

	histBaseURL := params.HistoricBaseURL
	if histBaseURL == nil {
		histBaseURL, err = url.Parse(defaultHistoricWhoisURL)
		if err != nil {
			panic(err)
		}
	}

	httpClient := http.DefaultClient
	if params.HTTPClient != nil {
		httpClient = params.HTTPClient
	}

	client := &Client{
		client:    httpClient,
		userAgent: userAgent,
		apiKey:    apiKey,
	}

	client.HistoricService = &historicServiceOp{client: client, baseURL: histBaseURL}

	return client
}

// Client is a client for Whois XML API services.
type Client struct {
	client *http.Client

	userAgent string
	apiKey    string

	HistoricService
}

// Response is a response wrapper.
type Response struct {
	*http.Response
}

// NewRequest creates a basic API request
func (c *Client) NewRequest(method string, u *url.URL, body io.Reader) (*http.Request, error) {

	var err error
	var req *http.Request

	req, err = http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", mediaType)
	req.Header.Add("Accept", mediaType)
	req.Header.Add("User-Agent", c.userAgent)

	return req, nil
}

// Do sends an API request and returns the API response
func (c *Client) Do(ctx context.Context, req *http.Request, v io.Writer) (response *Response, err error) {

	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot execute request: %w", err)
	}

	defer func() {
		if rerr := resp.Body.Close(); err == nil && rerr != nil {
			err = fmt.Errorf("cannot close response: %w", rerr)
		}
	}()

	response = &Response{Response: resp}

	_, err = io.Copy(v, resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response: %w", err)
	}

	return response, err
}

// ErrorResponse is returned when response's status code is not 2xx
type ErrorResponse struct {
	Response *http.Response
	Message  string
}

func (e ErrorResponse) Error() string {
	if e.Message != "" {
		return "API failed with status code: " + strconv.Itoa(e.Response.StatusCode) + " (" + e.Message + ")"
	}
	return "API failed with status code: " + strconv.Itoa(e.Response.StatusCode)
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	var errorResponse = ErrorResponse{
		Response: r,
	}

	return errorResponse
}
