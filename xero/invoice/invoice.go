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
	"encoding/json"
	"encoding/xml"
	"errors"
	"log"

	"github.com/arduino/go-xero/xero"
)

const (
	path = "/api.xro/2.0/invoices"
)

const (
	// InvoiceLineAmountTypeEXCLUSIVE - Invoice lines are exclusive of tax (default)
	InvoiceLineAmountTypeEXCLUSIVE = "Exclusive"
	// InvoiceLineAmountTypeINCLUSIVE - Invoice lines are inclusive tax
	InvoiceLineAmountTypeINCLUSIVE = "Inclusive"
	// InvoiceLineAmountTypeNOTAX - Invoices lines have no tax
	InvoiceLineAmountTypeNOTAX = "NoTax"
)

// Address is the Contact Address model
type Address struct {
	XMLName      xml.Name `xml:"Address"`
	AddressType  string
	AddressLine1 string
	AddressLine2 string
	City         string
	Region       string
	PostalCode   string
	Country      string
}

// ContactType is the Invoice Contact content model
type ContactType struct {
	XMLName   xml.Name `xml:"Contact"`
	Name      string
	Addresses []Address `xml:"Addresses>Address"`
}

// LineItemObj is the LineItem content model
type LineItemObj struct {
	Description string
	Quantity    string
	UnitAmount  string
	AccountCode string
}

// LineItem is the Invoice LineItem model
type LineItem struct {
	LineItem LineItemObj
}

// Invoice is the Invoice model
type Invoice struct {
	XMLName         xml.Name `xml:"Invoice"`
	InvoiceID       string   `xml:",omitempty"`
	InvoiceNumber   string   `xml:",omitempty"`
	Type            string
	Contact         ContactType
	Date            string
	DueDate         string
	Status          string
	LineAmountTypes string
	LineItems       LineItem
	Reference       string
}

type response struct {
	ID           string `xml:"Id"`
	Status       string `xml:",omitempty"`
	ProviderName string
	DateTimeUTC  string
	Invoices     []Invoice `xml:"Invoices>Invoice"`
}

// New creates an invoice
func New(inv Invoice) (resp string, err error) {

	var invoiceSaved response
	var errorResponse xero.ApiException

	xmlString, marshalErr := xml.Marshal(inv)
	if marshalErr != nil {
		log.Printf("error: %#v\n", marshalErr)
		return "", marshalErr
	}

	resp, err = xero.PostRequest(path, string(xmlString))
	if err != nil {
		log.Printf("[xero invoice New] - error: %#v\n", err)
		return "", err
	}

	//log.Printf("\n\n[xero invoice New] - Invoice Saved XML: %s\n", string(resp))
	savedMarshalErr := xml.Unmarshal([]byte(resp), &invoiceSaved)
	if savedMarshalErr != nil {
		log.Printf("[xero invoice New] - Xml Unmarshal Error: %#v\n", savedMarshalErr)
		return "", savedMarshalErr
	}

	if invoiceSaved.Status != "OK" {
		apiMarshalErr := xml.Unmarshal([]byte(resp), &errorResponse)
		if apiMarshalErr != nil {
			return "", apiMarshalErr
		}
		log.Printf("[xero invoice New] - Xero Api Error in Response: %#v\n", errorResponse)
		return "", errors.New(errorResponse.Message)
	}

	//log.Printf("\n\n[xero invoice New] - Invoice Saved: %#v\n", invoiceSaved)
	response := map[string]string{
		"InvoiceID":     invoiceSaved.Invoices[0].InvoiceID,
		"InvoiceNumber": invoiceSaved.Invoices[0].InvoiceNumber}
	jsonResponse, _ := json.Marshal(response)
	return string(jsonResponse), nil
}

// Query gets the invoices list
func Query() (resp string, err error) {

	//var invoicesList response

	resp, err = xero.Request("GET", path)
	if err != nil {
		log.Printf("[xero invoice Query] - error: %#v\n", err)
		return "", err
	}

	// log.Printf("\n\n[xero invoice Query] - Invoices List XML: %s\n", string(resp))
	// savedMarshalErr := xml.Unmarshal([]byte(resp), &invoicesList)
	// if savedMarshalErr != nil {
	// 	log.Printf("[xero invoices Query] - Xml Unmarshal Error: %#v\n", savedMarshalErr)
	// 	return "", savedMarshalErr
	// }
	// //log.Printf("\n\n[xero invoice Query] - Invoices List: %#v\n", invoicesList)
	//
	// if invoicesList.Status != "OK" {
	// 	log.Printf("[xero invoices Query] - Cannot List invoices, Status: %s", invoicesList.Status)
	// 	return "", fmt.Errorf("Cannot List invoices, Status: %s", invoicesList.Status)
	// }
	//
	// jsonResponse, _ := json.Marshal(invoicesList.Invoices)
	return resp, nil
}
