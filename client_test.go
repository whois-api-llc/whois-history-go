package whoishistory

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

const (
	pathWhoisResponseOK       = "/whois/ok"
	pathWhoisResponseError    = "/whois/error"
	pathWhoisResponse500      = "/whois/500"
	pathWhoisResponsePartial1 = "/whois/partial"
	pathWhoisResponsePartial2 = "/whois/partial2"
)

const apiKey = "at_LoremIpsumDolorSitAmetConsect"

func whoisServer(resp, respErr string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var response string

		response = resp

		switch req.URL.Path {
		case pathWhoisResponseOK:
		case pathWhoisResponseError:
			response = respErr
		case pathWhoisResponse500:
			w.WriteHeader(500)
			response = respErr
		case pathWhoisResponsePartial1:
			response = response[:len(response)-10]
		case pathWhoisResponsePartial2:
			w.Header().Set("Content-Length", strconv.Itoa(len(response)))
			response = response[:len(response)-10]
		default:
			panic(req.URL.Path)
		}
		_, err := w.Write([]byte(response))
		if err != nil {
			panic(err)
		}
	}))

	return server
}

func newAPI(apiServer *httptest.Server, link string) *Client {

	apiURL, err := url.Parse(apiServer.URL)
	if err != nil {
		panic(err)
	}
	apiURL.Path = link

	params := ClientParams{
		HTTPClient:      apiServer.Client(),
		WhoisBaseURL:    apiURL,
		HistoricBaseURL: apiURL,
	}

	return NewClient(apiKey, params)
}

func TestAPI_HistoricPurchase(t *testing.T) {

	checkResult := func(res int) bool {
		return res != 0
	}

	checkResultArr := func(res []*WhoisRecord) bool {
		return res != nil
	}

	ctx := context.Background()

	const resp = `{"recordsCount":1,"records":[{"domainName":"test.test"}]}`

	const errResp = `{"code":123,"messages":"test error"}`

	server := whoisServer(resp, errResp)
	defer server.Close()

	type args struct {
		ctx     context.Context
		options string
	}
	tests := []struct {
		name    string
		path    string
		args    args
		want    bool
		wantErr string
	}{
		{
			name: "successfull request",
			path: pathWhoisResponseOK,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    true,
			wantErr: "",
		},
		{
			name: "non 200 status code",
			path: pathWhoisResponse500,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "API error: [123] test error",
		},
		{
			name: "partial response 1",
			path: pathWhoisResponsePartial1,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot parse response: unexpected EOF",
		},
		{
			name: "partial response 2",
			path: pathWhoisResponsePartial2,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "cannot read response: unexpected EOF",
		},
		{
			name: "could not process request",
			path: pathWhoisResponseError,
			args: args{
				ctx:     ctx,
				options: "whoisxmlapi.com",
			},
			want:    false,
			wantErr: "API error: [123] test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			api := newAPI(server, tt.path)

			got, _, err := api.Preview(tt.args.ctx, tt.args.options)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("Whois.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want {
				if !checkResult(got) {
					t.Errorf("Whois.Get() got = %v, expected something else", got)
				}
			} else {
				if got != 0 {
					t.Errorf("Whois.Get() got = %v, expected nil", got)
				}
			}

			gotArr, _, err := api.Purchase(tt.args.ctx, tt.args.options)
			if (err != nil || tt.wantErr != "") && (err == nil || err.Error() != tt.wantErr) {
				t.Errorf("Whois.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.want {
				if !checkResultArr(gotArr) {
					t.Errorf("Whois.Get() got = %v, expected something else", gotArr)
				}
			} else {
				if gotArr != nil {
					t.Errorf("Whois.Get() got = %v, expected nil", gotArr)
				}
			}
		})
	}
}
