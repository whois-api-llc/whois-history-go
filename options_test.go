package whoishistory

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

func TestOptions(t *testing.T) {

	d := time.Date(2020, 01, 01, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name   string
		values url.Values
		option Option
		want   string
	}{
		{
			name:   "since date",
			values: url.Values{},
			option: OptionSinceDate(d),
			want:   "sinceDate=2020-01-01",
		},
		{
			name:   "created from date",
			values: url.Values{},
			option: OptionCreatedDateFrom(d),
			want:   "createdDateFrom=2020-01-01",
		},
		{
			name:   "created to date",
			values: url.Values{},
			option: OptionCreatedDateTo(d),
			want:   "createdDateTo=2020-01-01",
		},
		{
			name:   "updated from date",
			values: url.Values{},
			option: OptionUpdatedDateFrom(d),
			want:   "updatedDateFrom=2020-01-01",
		},
		{
			name:   "updated to date",
			values: url.Values{},
			option: OptionUpdatedDateTo(d),
			want:   "updatedDateTo=2020-01-01",
		},
		{
			name:   "expired from date",
			values: url.Values{},
			option: OptionExpiredDateFrom(d),
			want:   "expiredDateFrom=2020-01-01",
		},
		{
			name:   "expired to date",
			values: url.Values{},
			option: OptionExpiredDateTo(d),
			want:   "expiredDateTo=2020-01-01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.option(tt.values)
			if got := tt.values.Encode(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Option() = %v, want %v", got, tt.want)
			}
		})
	}
}
