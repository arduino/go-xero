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
	"crypto/x509"
	"encoding/pem"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/garyburd/go-oauth/oauth"
	jww "github.com/spf13/jwalterweatherman"
)

const (
	requestTokenURL   = "https://api.xero.com/oauth/RequestToken"
	authorizeTokenURL = "https://api.xero.com/oauth/Authorize"
	accessTokenURL    = "https://api.xero.com/oauth/AccessToken"
	baseURL           = "https://api.xero.com"
)

// this is global, it's very bad!
// var client oauth.Client

// APIException model for error response
type APIException struct {
	XMLName     xml.Name `xml:"ApiException"`
	Type        string
	ErrorNumber int
	Message     string
}

// Xoauth is a wrapper around oauth.Client
type Xoauth struct {
	*oauth.Client
}

// NewClient initializes the oauth Client
func NewClient(token string, key []byte) (client Xoauth, err error) {

	block, _ := pem.Decode(key)
	privateKey, ParseKeyErr := x509.ParsePKCS1PrivateKey(block.Bytes)

	if ParseKeyErr != nil {
		jww.ERROR.Printf("[xero NewClient] - Parse private key ERROR: %#v", ParseKeyErr)
		return Xoauth{}, ParseKeyErr
	}

	client = Xoauth{
		&oauth.Client{
			TemporaryCredentialRequestURI: requestTokenURL,
			ResourceOwnerAuthorizationURI: authorizeTokenURL,
			TokenRequestURI:               accessTokenURL,
			//Header:                        http.Header{"Accept": {"application/json"}},
			SignatureMethod: oauth.RSASHA1,
			Credentials:     oauth.Credentials{Token: token},
			PrivateKey:      privateKey,
		}}

	// myclient := Xoauth{&client}
	// myclient := &client
	return client, ParseKeyErr
}

// PostRequest sends POST requests to xero APIs with a form as payload
func (client Xoauth) PostRequest(path string, payload string) (response string, err error) {

	form := url.Values{"xml": {payload}}

	req, reqErr := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if reqErr != nil {
		jww.ERROR.Printf("[xero PostRequest] - Error: %#v\n", reqErr)
		return "", reqErr
	}

	req.URL.Path = path

	headerErr := client.SetAuthorizationHeader(req.Header, &client.Credentials, "POST", req.URL, nil)
	if headerErr != nil {
		jww.ERROR.Printf("[xero PostRequest] - SetAuthorizationHeader Error: %#v\n", headerErr)
		return "", headerErr
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.Header.Set("Accept", "application/json")

	resp, reqErr := client.Post(http.DefaultClient, &client.Credentials, req.URL.String(), form)
	if reqErr != nil {
		jww.ERROR.Printf("[xero PostRequest] - Error: %#v\n", reqErr)
		return "", reqErr
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return string(body), nil
}

// Options adds optional parameters to Xero requests
type Options struct {
	ModifiedAfter string
	Values        url.Values
}

// Request sends requests to xero APIs
func (client Xoauth) Request(method string, path string, otherOptions *Options) (response string, err error) {

	req, err := http.NewRequest(method, baseURL, nil)
	jww.DEBUG.Printf("[xero Request] -  req in NewRequest: %#v\n", req.URL.String())
	if err != nil {
		jww.ERROR.Printf("[xero Request] - error in NewRequest: %#v\n", err)
		return "", err
	}

	if otherOptions != nil {
		req.URL.RawQuery = otherOptions.Values.Encode()
	}

	req.URL.Path = path

	headerErr := client.SetAuthorizationHeader(req.Header, &client.Credentials, method, req.URL, nil)
	if headerErr != nil {
		jww.ERROR.Printf("[xero Request] - SetAuthorizationHeader Error: %#v\n", headerErr)
		return "", headerErr
	}
	// jww.DEBUG.Printf("other options: %v", otherOptions)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	// req.Header.Set("Accept", "application/json")
	// jww.DEBUG.Printf("[xero Request] - req: %v\n", req)
	resp, reqErr := http.DefaultClient.Do(req)
	if reqErr != nil {
		jww.ERROR.Printf("[xero Request] - Error in Do: %#v %#v\n", req, reqErr)
		return "", reqErr
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return string(body), nil

}
