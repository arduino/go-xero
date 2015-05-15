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
	"encoding/xml"
	"errors"

	"github.com/arduino/go-xero/xero/payment"
	jww "github.com/spf13/jwalterweatherman"
)

// NewPayment creates one or more payments for the given Invoices
func (client Xoauth) NewPayment(paym []payment.Payment) ([]string, error) {

	var paymentsToSave payment.Payments

	paymentsToSave.Payments = paym

	var paymentSaved payment.Response
	var errorResponse APIException

	xmlString, marshalErr := xml.Marshal(paymentsToSave)
	if marshalErr != nil {
		jww.ERROR.Printf("error: %#v\n", marshalErr)
		return nil, marshalErr
	}

	resp, err := client.PostRequest(payment.Path, string(xmlString))
	if err != nil {
		jww.ERROR.Printf("[xero payment New] - error: %#v\n", err)
		return nil, err
	}

	//jww.DEBUG.Printf("\n\n[xero payment New] - Payment Saved XML: %s\n", string(resp))
	savedMarshalErr := xml.Unmarshal([]byte(resp), &paymentSaved)
	if savedMarshalErr != nil {
		jww.ERROR.Printf("[xero payment New] - Xml Unmarshal Error: %#v\n", savedMarshalErr)
		return nil, savedMarshalErr
	}

	if paymentSaved.Status != "OK" {
		apiMarshalErr := xml.Unmarshal([]byte(resp), &errorResponse)
		if apiMarshalErr != nil {
			return nil, apiMarshalErr
		}
		jww.ERROR.Printf("[xero payment New] - Xero Api Error in Response: %#v\n", errorResponse)
		return nil, errors.New(errorResponse.Message)
	}

	var payments []string
	for _, payment := range paymentSaved.Payments {
		payments = append(payments, payment.PaymentID)
	}

	return payments, nil

}
