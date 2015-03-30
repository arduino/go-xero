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

package invoice

import (
	"encoding/xml"
	"log"

	"github.com/arduino/go-xero/xero"
)

const (
	path = "/api.xro/2.0/invoices"
)

type address struct {
	XMLName     xml.Name `xml:"address"`
	AddressType string
	Country     string
}

type contactType struct {
	XMLName   xml.Name `xml:"Contact"`
	Name      string
	Addresses []address
}

type lineItemObj struct {
	Description string
	Quantity    string
	UnitAmount  string
	AccountCode string
}

type lineItem struct {
	LineItem lineItemObj
}

// Invoice is the invoice model
type Invoice struct {
	XMLName         xml.Name `xml:"Invoice"`
	Type            string
	Contact         contactType
	Date            string
	DueDate         string
	LineAmountTypes string
	LineItems       lineItem
}

// New creates an invoice
func New(inv Invoice) (resp string, err error) {

	xmlString, marshalErr := xml.Marshal(inv)
	if marshalErr != nil {
		log.Printf("error: %v\n", marshalErr)
	}

	resp, err = xero.PostRequest(path, string(xmlString))
	if err != nil {
		log.Printf("[xero invoice New] - error: %v\n", err)
		return "", err
	}
	return resp, nil
}

// Query gets the invoices list
func Query() (resp string, err error) {

	resp, err = xero.Request("GET", path)
	if err != nil {
		log.Printf("[xero invoice Query] - error: %v\n", err)
		return "", err
	}

	return resp, nil
}
