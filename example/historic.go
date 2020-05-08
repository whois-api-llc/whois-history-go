package example

import (
	"context"
	"errors"
	"github.com/whois-api-llc/whois-history-go"
	"log"
	"time"
)

// HistoricWhoisPreview is an example of historic whois API usage
func HistoricWhoisPreview(apiKey string) {

	client := whoishistory.NewBasicClient(apiKey)

	// Create request with additional options. Options can increase response time.
	num, _, err := client.Preview(context.Background(), "whoisxmlapi.com",
		// Let's check records registered in 2019
		whoishistory.OptionCreatedDateFrom(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)),
		whoishistory.OptionCreatedDateTo(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
	)
	if err != nil {
		// Handle error message returned by server
		var apiErr *whoishistory.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	log.Println(num)
}

// HistoricWhoisPurchase is an example of historic whois API usage
func HistoricWhoisPurchase(apiKey string) {

	client := whoishistory.NewBasicClient(apiKey)

	// Create request with additional options. Options can increase response time.
	num, _, err := client.HistoricService.Preview(context.Background(), "whoisxmlapi.com",
		// Let's check records registered in 2019
		whoishistory.OptionCreatedDateFrom(time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)),
		whoishistory.OptionCreatedDateTo(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)),
	)
	if err != nil {
		// Handle error message returned by server
		var apiErr *whoishistory.ErrorMessage
		if errors.As(err, &apiErr) {
			log.Println(apiErr.Code)
			log.Println(apiErr.Message)
		}
		log.Fatal(err)
	}

	log.Println(num)
}
