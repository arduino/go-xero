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

package xero

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"net/url"
	"strconv"
	"time"

	"github.com/arduino/go-xero/xero/invoice"
	jww "github.com/spf13/jwalterweatherman"
)

// GetAllInvoices gives you all the invoices of the org
func (client Xoauth) GetAllInvoices() (allInvoices invoice.Invoices, err error) {

	var invoiceOptions Options
	var invoiceList invoice.Response
	var responseList []invoice.Response

	// invoiceOptions.ModifiedAfter = "2014-01-01T00:00:00"
	invoiceOptions.Values = url.Values{}
	// this is a do while statement
	// it stops paging if the are no more invoices
	for i := 1; ; i++ {
		invoiceOptions.Values.Set("page", strconv.Itoa(i))

		response, reqErr := client.Request("GET", invoice.Path, &invoiceOptions)
		if reqErr != nil {
			jww.ERROR.Printf("[xero invoice GetAllInvoices] - Error response: %#v", reqErr)
		}
		invoicesMarshalErr := xml.Unmarshal([]byte(response), &invoiceList)
		if invoicesMarshalErr != nil {
			jww.ERROR.Printf("[xero invoice GetAllInvoices] - Xml Unmarshal Error: %#v\n", invoicesMarshalErr)
			return allInvoices, invoicesMarshalErr
		}
		// jww.DEBUG.Printf("invoice list: %v", invoiceList.Invoices)
		responseList = append(responseList, invoiceList)

		// if there are no more invoices to fecth, stop asking for them
		if len(invoiceList.Invoices) <= 0 {
			break
		}
		// clean up the invoice list for the next request
		invoiceList = invoice.Response{}
		// avoid xero limit, there is a limit of 60 req/min we wait 60s every 50 reqs
		if i%50 == 0 {
			jww.DEBUG.Printf("i:%d", i)
			time.Sleep(60 * time.Second)
		}
	}

	for singleResponse := range responseList {
		allInvoices.Invoices = append(allInvoices.Invoices, responseList[singleResponse].Invoices...)
	}
	// jww.DEBUG.Printf("tipo b: %v", allInvoices.Invoices)
	return allInvoices, nil
}

// NewInvoice creates one or more invoices
func (client Xoauth) NewInvoice(inv []invoice.Invoice) (resp string, err error) {

	var invoiceSaved invoice.Response
	var errorResponse APIException

	var invoicesToSave invoice.Invoices

	invoicesToSave.Invoices = inv
	xmlString, marshalErr := xml.Marshal(invoicesToSave)
	if marshalErr != nil {
		jww.ERROR.Printf("error: %#v\n", marshalErr)
		return "", marshalErr
	}
	//jww.ERROR.Printf("\n\n[xero invoice New] - Invoice XML to send: %s\n", string(xmlString))
	resp, err = client.PostRequest(invoice.Path, string(xmlString))
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

	jww.ERROR.Printf("\n\n[xero invoice New] - Invoice Saved: %#v\n", invoiceSaved)
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
