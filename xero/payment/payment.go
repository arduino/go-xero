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

import (
	"encoding/xml"
	"errors"
	"log"

	"github.com/arduino/go-xero/xero"
)

const (
	path = "/api.xro/2.0/Payments"
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

type response struct {
	ID           string `xml:"Id"`
	Status       string `xml:",omitempty"`
	ProviderName string
	DateTimeUTC  string
	Payments     []Payment `xml:"Payments>Payment"`
}

// New creates a payment for the given Invoice
func New(paym Payment) (resp string, err error) {

	var paymentSaved response
	var errorResponse xero.ApiException

	xmlString, marshalErr := xml.Marshal(paym)
	if marshalErr != nil {
		log.Printf("error: %#v\n", marshalErr)
		return "", marshalErr
	}

	resp, err = xero.PostRequest(path, string(xmlString))
	if err != nil {
		log.Printf("[xero payment New] - error: %#v\n", err)
		return "", err
	}

	//log.Printf("\n\n[xero payment New] - Payment Saved XML: %s\n", string(resp))
	savedMarshalErr := xml.Unmarshal([]byte(resp), &paymentSaved)
	if savedMarshalErr != nil {
		log.Printf("[xero payment New] - Xml Unmarshal Error: %#v\n", savedMarshalErr)
		return "", savedMarshalErr
	}

	if paymentSaved.Status != "OK" {
		apiMarshalErr := xml.Unmarshal([]byte(resp), &errorResponse)
		if apiMarshalErr != nil {
			return "", apiMarshalErr
		}
		log.Printf("[xero payment New] - Xero Api Error in Response: %#v\n", errorResponse)
		return "", errors.New(errorResponse.Message)
	}

	return paymentSaved.Payments[0].PaymentID, nil

}
