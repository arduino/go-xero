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

import "encoding/xml"

const (
	// Path is the relative API path for invoices
	Path = "/api.xro/2.0/Invoices"
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

// Invoices is the Invoices array model
type Invoices struct {
	Invoices []Invoice `xml:"Invoices"`
}

// Response is the xero request response model for invoices
type Response struct {
	ID           string `xml:"Id"`
	Status       string `xml:",omitempty"`
	ProviderName string
	DateTimeUTC  string
	Invoices     []Invoice `xml:"Invoices>Invoice"`
}
