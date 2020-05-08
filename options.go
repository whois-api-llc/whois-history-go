package whoishistory

import (
	"net/url"
	"time"
)

// Option adds parameters to the query
type Option func(v url.Values)

var _ = []Option{
	OptionSinceDate(time.Time{}),
	OptionCreatedDateFrom(time.Time{}),
	OptionCreatedDateTo(time.Time{}),
	OptionUpdatedDateFrom(time.Time{}),
	OptionUpdatedDateTo(time.Time{}),
	OptionExpiredDateFrom(time.Time{}),
	OptionExpiredDateTo(time.Time{}),
}

const dateFormat = "2006-01-02"

// OptionSinceDate filters activities discovered since the given date.
// Sometimes there is a latency between the actual added/renewal/expired date
// and the date when our system detected this change.
// This function ignores time and sends only a date in format 2006-01-02
func OptionSinceDate(date time.Time) Option {
	return func(v url.Values) {
		v.Set("sinceDate", date.Format(dateFormat))
	}
}

// OptionCreatedDateFrom searches through domains created after the given date.
func OptionCreatedDateFrom(date time.Time) Option {
	return func(v url.Values) {
		v.Set("createdDateFrom", date.Format(dateFormat))
	}
}

// OptionCreatedDateTo searches through domains created before the given date.
func OptionCreatedDateTo(date time.Time) Option {
	return func(v url.Values) {
		v.Set("createdDateTo", date.Format(dateFormat))
	}
}

// OptionUpdatedDateFrom searches through domains updated after the given date.
func OptionUpdatedDateFrom(date time.Time) Option {
	return func(v url.Values) {
		v.Set("updatedDateFrom", date.Format(dateFormat))
	}
}

// OptionUpdatedDateTo searches through domains updated before the given date.
func OptionUpdatedDateTo(date time.Time) Option {
	return func(v url.Values) {
		v.Set("updatedDateTo", date.Format(dateFormat))
	}
}

// OptionExpiredDateFrom searches through domains expired after the given date.
func OptionExpiredDateFrom(date time.Time) Option {
	return func(v url.Values) {
		v.Set("expiredDateFrom", date.Format(dateFormat))
	}
}

// OptionExpiredDateTo searches through domains expired before the given date.
func OptionExpiredDateTo(date time.Time) Option {
	return func(v url.Values) {
		v.Set("expiredDateTo", date.Format(dateFormat))
	}
}
