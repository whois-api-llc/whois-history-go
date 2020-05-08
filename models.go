package whoishistory

import (
	"encoding/json"
	"fmt"
	"time"
)

func unmarshalString(raw json.RawMessage) (string, error) {
	var val string
	err := json.Unmarshal(raw, &val)
	if err != nil {
		return "", err
	}
	return val, nil
}

// Time is a helper wrapper on time.Time
type Time time.Time

var emptyTime Time

// UnmarshalJSON decodes time as historic whois API does
func (t *Time) UnmarshalJSON(b []byte) error {
	str, err := unmarshalString(b)
	if err != nil {
		return err
	}
	if str == "" {
		*t = emptyTime
		return nil
	}
	v, err := time.Parse("2006-01-02T15:04:05-07:00", str)
	if err != nil {
		return err
	}
	*t = Time(v)
	return nil
}

// MarshalJSON encodes time as historic whois API does
func (t Time) MarshalJSON() ([]byte, error) {
	if t == emptyTime {
		return []byte(`""`), nil
	}
	return []byte(`"` + time.Time(t).Format("2006-01-02T15:04:05-07:00") + `"`), nil
}

// Audit is a part of whois API response. It represents dates
// when whois record was added and updated in our database.
type Audit struct {
	CreatedDate Time `json:"createdDate"`
	UpdatedDate Time `json:"updatedDate"`
}

// Contact is a part of historic whois API response
type Contact struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
	Street       string `json:"street"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postalCode"`
	Country      string `json:"country"`
	Email        string `json:"email"`
	Telephone    string `json:"telephone"`
	TelephoneExt string `json:"telephoneExt"`
	Fax          string `json:"fax"`
	FaxExt       string `json:"faxExt"`
	RawText      string `json:"rawText"`
}

// WhoisRecord is a whois record
type WhoisRecord struct {
	DomainName            string   `json:"domainName"`
	DomainType            string   `json:"domainType"`
	CreatedDateISO8601    Time     `json:"createdDateISO8601"`
	UpdatedDateISO8601    Time     `json:"updatedDateISO8601"`
	ExpiresDateISO8601    Time     `json:"expiresDateISO8601"`
	CreatedDateRaw        string   `json:"createdDateRaw"`
	UpdatedDateRaw        string   `json:"updatedDateRaw"`
	ExpiresDateRaw        string   `json:"expiresDateRaw"`
	Audit                 Audit    `json:"audit"`
	NameServers           []string `json:"nameServers"`
	WhoisServer           string   `json:"whoisServer"`
	RegistrarName         string   `json:"registrarName"`
	Status                []string `json:"status"`
	CleanText             string   `json:"cleanText"`
	RawText               string   `json:"rawText"`
	RegistrantContact     Contact  `json:"registrantContact"`
	AdministrativeContact Contact  `json:"administrativeContact"`
	TechnicalContact      Contact  `json:"technicalContact"`
	BillingContact        Contact  `json:"billingContact"`
	ZoneContact           Contact  `json:"zoneContact"`
}

// ErrorMessage is a error message from historic whois API
type ErrorMessage struct {
	Code    int    `json:"code"`
	Message string `json:"messages"`
}

func (e ErrorMessage) Error() string {
	return fmt.Sprintf("API error: [%d] %s", e.Code, e.Message)
}
