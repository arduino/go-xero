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
	"net/url"
	"strconv"
	"time"

	"github.com/arduino/go-xero/xero"
	jww "github.com/spf13/jwalterweatherman"
)

const (
	path = "/api.xro/2.0/Invoices"
)

const (
	// InvoiceLineAmountTypeEXCLUSIVE - Invoice lines are exclusive of tax (default)
	InvoiceLineAmountTypeEXCLUSIVE = "Exclusive"
	// InvoiceLineAmountTypeINCLUSIVE - Invoice lines are inclusive tax
	InvoiceLineAmountTypeINCLUSIVE = "Inclusive"
	// InvoiceLineAmountTypeNOTAX - Invoices lines have no tax
	InvoiceLineAmountTypeNOTAX = "NoTax"
)

type Xclient struct {
	client xero.Xoauth
}

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
	XMLName             xml.Name `xml:"Invoice"`
	InvoiceID           string   `xml:",omitempty"`
	InvoiceNumber       string   `xml:",omitempty"`
	Type                string
	Contact             ContactType
	Date                string
	DueDate             string
	ExpectedPaymentDate string
	Status              string
	LineAmountTypes     string
	LineItems           LineItem
	Reference           string
}

type invoices struct {
	Invoices []Invoice `xml:"Invoices"`
}

type response struct {
	ID           string `xml:"Id"`
	Status       string `xml:",omitempty"`
	ProviderName string
	DateTimeUTC  string
	Invoices     []Invoice `xml:"Invoices>Invoice"`
}

// New creates one or more invoices
func (invClient Xclient) New(inv []Invoice) (resp string, err error) {

	var invoiceSaved response
	var errorResponse xero.APIException

	var invoicesToSave invoices

	invoicesToSave.Invoices = inv
	xmlString, marshalErr := xml.Marshal(invoicesToSave)
	if marshalErr != nil {
		jww.ERROR.Printf("error: %#v\n", marshalErr)
		return "", marshalErr
	}
	//jww.DEBUG.Printf("\n\n[xero invoice New] - Invoice XML to send: %s\n", string(xmlString))
	resp, err = invClient.client.PostRequest(path, string(xmlString))
	if err != nil {
		jww.ERROR.Printf("[xero invoice New] - error: %#v\n", err)
		return "", err
	}

	//jww.DEBUG.Printf("\n\n[xero invoice New] - Invoice Saved XML: %s\n", string(resp))
	savedMarshalErr := xml.Unmarshal([]byte(resp), &invoiceSaved)
	if savedMarshalErr != nil {
		jww.ERROR.Printf("[xero invoice New] - Xml Unmarshal Error: %#v\n", savedMarshalErr)
		return "", savedMarshalErr
	}

	if invoiceSaved.Status != "OK" {
		apiMarshalErr := xml.Unmarshal([]byte(resp), &errorResponse)
		if apiMarshalErr != nil {
			return "", apiMarshalErr
		}
		jww.ERROR.Printf("[xero invoice New] - Xero Api Error in Response: %#v\n", errorResponse)
		return "", errors.New(errorResponse.Message)
	}

	jww.DEBUG.Printf("\n\n[xero invoice New] - Invoice Saved: %#v\n", invoiceSaved)
	var itemsSaved []map[string]string
	for _, invoice := range invoiceSaved.Invoices {
		item := map[string]string{
			"InvoiceID":           invoice.InvoiceID,
			"InvoiceNumber":       invoice.InvoiceNumber,
			"Reference":           invoice.Reference,
			"ExpectedPaymentDate": invoice.ExpectedPaymentDate,
			"Amount":              invoice.LineItems.LineItem.UnitAmount}
		itemsSaved = append(itemsSaved, item)
	}

	jsonResponse, _ := json.Marshal(itemsSaved)
	return string(jsonResponse), nil
}

// GetAllInvoices gives you all the invoices of the org
func (invClient Xclient) GetAllInvoices() (invoices []response, err error) {

	var invoiceOptions xero.Options
	var invoiceList response
	var responseList []response

	// invoiceOptions.ModifiedAfter = "2014-01-01T00:00:00"
	invoiceOptions.Values = url.Values{}
	for i := 1; i <= 30; i++ {
		invoiceOptions.Values.Set("page", strconv.Itoa(i))
		resp, err1 := invClient.client.Request("GET", "/api.xro/2.0/Invoices", &invoiceOptions)
		if err1 != nil {
			jww.ERROR.Printf("[xero invoice GetAllInvoices] - Error response: %#v\n %#v", resp, err1)
		}
		invoicesMarshalErr := xml.Unmarshal([]byte(resp), &invoiceList)
		if invoicesMarshalErr != nil {
			jww.ERROR.Printf("[xero invoice GetAllInvoices] - Xml Unmarshal Error: %#v\n", invoicesMarshalErr)
			return invoices, invoicesMarshalErr
		}
		responseList = append(responseList, invoiceList)
		// clean up the invoice list for the next request
		invoiceList = response{}

		if i%50 == 0 {
			jww.DEBUG.Printf("i:%d", i)
			time.Sleep(60 * time.Second)
		}
	}

	// jww.DEBUG.Printf("responseList: %#v\n", pretty.Formatter(responseList))

	return responseList, nil

}
