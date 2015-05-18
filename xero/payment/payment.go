/*
  Copyright 2015 Arduino LLC (http://www.arduino.cc/)

  This file is part of go-xero.

  go-xero is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, version 3 of the License,
  go-xero is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with go-xero.  If not, see <http://www.gnu.org/licenses/>.
*/

package payment

import "encoding/xml"

const (
	// Path is the relative API path for Payments
	Path = "/api.xro/2.0/Payments"
)

// InvoiceParam is the model with the Invoice Number
type InvoiceParam struct {
	InvoiceNumber string
}

// AccountParam is the model with the Account Code
type AccountParam struct {
	Code string
}

// Payment is the Payment model
type Payment struct {
	XMLName   xml.Name `xml:"Payment"`
	PaymentID string   `xml:",omitempty"`
	Invoice   InvoiceParam
	Account   AccountParam
	Date      string
	Amount    string
	Reference string
}

// Payments is the Payments array model
type Payments struct {
	Payments []Payment
}

// Response is the xero request response model for payments
type Response struct {
	ID           string `xml:"Id"`
	Status       string `xml:",omitempty"`
	ProviderName string
	DateTimeUTC  string
	Payments     []Payment `xml:"Payments>Payment"`
}
