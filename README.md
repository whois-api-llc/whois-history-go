# Overview

The client library for
[Whois History API](https://whois-history.whoisxmlapi.com/)
in Go language.

The minimum go version is 1.13.

# Installation

The library is distributed as a Go module

```bash
go get github.com/whois-api-llc/whois-history-go
```

# Examples

Full API documentation available [here](https://whois-history.whoisxmlapi.com/api/documentation/making-requests)

You can find all examples in `example` directory.

## Create a new client

To start making requests you need the API Key. 
You can find it on your profile page on [whoisxmlapi.com](https://whoisxmlapi.com/).
Using the API Key you can create Client.

Most users will be fine with `NewBasicClient` function. 
```go
client := whoisxmlapigo.NewBasicClient(apikey)
```

If you want to set custom `http.Client` to use proxy then you can use `NewClient` function.
```go
transport := &http.Transport{Proxy: http.ProxyURL(proxyUrl)}

client := whoisxmlapigo.NewClient(apiKey, whoisxmlapigo.ClientParams{
    HTTPClient: &http.Client{
        Transport: transport,
        Timeout:   20 * time.Second,
    },
})
```

## Make basic requests

Whois History API provides the historic registration details of a domain name. 

```go

// Make request to check number of records available for domain
num, _, err := client.HistoricService.Preview(ctx, "whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

log.Println(num)

// Make request to get actual records for domain
records, _, err := client.HistoricService.Purchase(ctx, "whoisxmlapi.com")
if err != nil {
    log.Fatal(err)
}

for _, rec := range records {
    log.Println(rec.Audit.UpdatedDate, rec.RegistrarName)
}
```
