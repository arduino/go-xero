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

package manualjournal

import (
	"encoding/json"
	"encoding/xml"
	"errors"

	"github.com/arduino/go-xero/xero"
	jww "github.com/spf13/jwalterweatherman"
)

const (
	path = "/api.xro/2.0/ManualJournals"
)

type Xclient struct {
	client xero.Xoauth
}

// JournalLineObj is the JournalLine single object model
type JournalLineObj struct {
	LineAmount  string
	AccountCode string
}

// JournalLine is the JournalLine model
type JournalLine struct {
	JournalLine JournalLineObj
}

// Journal is the Manual Journal model
type Journal struct {
	XMLName         xml.Name `xml:"ManualJournal"`
	ManualJournalID string   `xml:",omitempty"`
	Date            string
	Status          string
	Narration       string
	JournalLines    []JournalLine
}

type response struct {
	ID           string `xml:"Id"`
	Status       string `xml:",omitempty"`
	ProviderName string
	DateTimeUTC  string
	Journals     []Journal `xml:"ManualJournals>ManualJournal"`
}

// New creates one Journal
func (journClient Xclient) New(journ Journal) (resp string, err error) {

	var journalSaved response
	var errorResponse xero.APIException

	xmlString, marshalErr := xml.Marshal(journ)
	if marshalErr != nil {
		jww.ERROR.Printf("error: %#v\n", marshalErr)
		return "", marshalErr
	}
	//jww.ERROR.Printf("\n\n[xero manual jorunal New] - Journal XML to send: %s\n", string(xmlString))
	resp, err = journClient.client.PostRequest(path, string(xmlString))
	if err != nil {
		jww.ERROR.Printf("[xero manual journal New] - error: %#v\n", err)
		return "", err
	}

	//jww.ERROR.Printf("\n\n[xero manual journal New] - Manual Journal Saved XML: %s\n", string(resp))
	savedMarshalErr := xml.Unmarshal([]byte(resp), &journalSaved)
	if savedMarshalErr != nil {
		jww.ERROR.Printf("[xero manual journal New] - Xml Unmarshal Error: %#v\n", savedMarshalErr)
		return "", savedMarshalErr
	}

	if journalSaved.Status != "OK" {
		apiMarshalErr := xml.Unmarshal([]byte(resp), &errorResponse)
		if apiMarshalErr != nil {
			return "", apiMarshalErr
		}
		jww.ERROR.Printf("[xero manual journal New] - Xero Api Error in Response: %#v\n", errorResponse)
		return "", errors.New(errorResponse.Message)
	}

	jww.DEBUG.Printf("\n\n[xero manual journal New] - Manual Journal Saved: %#v\n", journalSaved)
	var itemsSaved []map[string]string
	for _, journal := range journalSaved.Journals {
		item := map[string]string{
			"JournalID": journal.ManualJournalID,
			"Date":      journal.Date}
		itemsSaved = append(itemsSaved, item)
	}

	jsonResponse, _ := json.Marshal(itemsSaved)
	return string(jsonResponse), nil
}
