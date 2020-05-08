package whoishistory

import (
	"encoding/json"
	"testing"
)

func TestTime(t *testing.T) {
	tests := []struct {
		name   string
		decErr string
		encErr string
	}{
		{
			name:   `"2006-01-02T15:04:05-07:00"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"2006-01-02T15:04:05-08:00"`,
			decErr: "",
			encErr: "",
		},
		{
			name:   `"2006-01-02T15:04:05Z08:00"`,
			decErr: `parsing time "2006-01-02T15:04:05Z08:00" as "2006-01-02T15:04:05-07:00": cannot parse "" as "-07:00"`,
			encErr: "",
		},
		{
			name:   `""`,
			decErr: "",
			encErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var v Time

			err := json.Unmarshal([]byte(tt.name), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.name {
				t.Errorf("got = %v, want %v", string(bb), tt.name)
			}
		})
	}
}

func TestContact(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
		decErr string
		encErr string
	}{
		{
			name:   `test-1`,
			input:  `{}`,
			output: `{"name":"","organization":"","street":"","city":"","state":"","postalCode":"","country":"","email":"","telephone":"","telephoneExt":"","fax":"","faxExt":"","rawText":""}`,
			decErr: "",
			encErr: "",
		},
		{
			name: `test-2`,
			input: `{
        "name": "cont-name",
        "organization": "cont-org",
        "street": "cont-street",
        "city": "cont-city",
        "state": "cont-state",
        "postalCode": "cont-postalCode",
        "country": "cont-country",
        "email": "cont-email",
        "telephone": "cont-telephone",
        "telephoneExt": "cont-telephoneExt",
        "fax": "cont-fax",
        "faxExt": "cont-faxExt",
        "rawText": "cont-rawText"
      }`,
			output: `{"name":"cont-name","organization":"cont-org","street":"cont-street","city":"cont-city","state":"cont-state","postalCode":"cont-postalCode","country":"cont-country","email":"cont-email","telephone":"cont-telephone","telephoneExt":"cont-telephoneExt","fax":"cont-fax","faxExt":"cont-faxExt","rawText":"cont-rawText"}`,
			decErr: "",
			encErr: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var v Contact

			err := json.Unmarshal([]byte(tt.input), &v)
			checkErr(t, err, tt.decErr)
			if tt.decErr != "" {
				return
			}

			bb, err := json.Marshal(v)
			checkErr(t, err, tt.encErr)
			if tt.encErr != "" {
				return
			}

			if string(bb) != tt.output {
				t.Errorf("got  = %v", string(bb))
				t.Errorf("want = %v", tt.output)
			}
		})
	}
}

func checkErr(t *testing.T, err error, want string) {
	if (err != nil || want != "") && (err == nil || err.Error() != want) {
		t.Errorf("error = %v, wantErr %v", err, want)
	}
}
