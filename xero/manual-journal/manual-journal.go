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

import "encoding/xml"

const (
	// Path is the relative API path for Manual Journals
	Path = "/api.xro/2.0/ManualJournals"
)

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

// Response is the xero request response model for journals
type Response struct {
	ID           string `xml:"Id"`
	Status       string `xml:",omitempty"`
	ProviderName string
	DateTimeUTC  string
	Journals     []Journal `xml:"ManualJournals>ManualJournal"`
}
