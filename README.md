#Go Xero

Implements interactions with a Xero private application.

## Installation

You need to have Git and Go already installed.
Run this in your terminal

```sh
go get github.com/arduino/go-xero
```

## Usage

Import it in your Go code:

```go
import (
  "github.com/arduino/go-xero/xero"
  "github.com/arduino/go-xero/xero/invoice"
)
```

## Client Creation

To initialize a client you need a private key and a consumer key

```go
keyFile, openFileErr := ioutil.ReadFile("privatekey.pem")
if openFileErr != nil {
  log.Fatal("opening key file ERROR: ", openFileErr)
  return
}

xeroClientErr := xero.NewClient("your_consumer_key", keyFile)
if xeroClientErr != nil {
  log.Fatal("init xero client ERROR: ", xeroClientErr)
  return
}
```

Get the invoices list

```go
invoice.Query()
```

Create new invoice

```go
var invoiceToSave invoice.Invoice
// ... invoiceToSave marshal / population
r, err := invoice.New(invoiceToSave)
```
