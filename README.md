#Go Xero

## Installation

```sh
go get github.com/arduino/go-xero
```

### Invoices

```go
import (
  "github.com/arduino/go-xero/xero"
  "github.com/arduino/go-xero/xero/invoice"
)

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

// get Invoices
r, err := invoice.Query()

// create Invoice
var invoiceToSave invoice.Invoice

// ... invoiceToSave marshal / population

r, err := invoice.New(invoiceToSave)
```
