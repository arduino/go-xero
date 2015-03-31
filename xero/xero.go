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
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/garyburd/go-oauth/oauth"
)

const (
	requestTokenURL   = "https://api.xero.com/oauth/RequestToken"
	authorizeTokenURL = "https://api.xero.com/oauth/Authorize"
	accessTokenURL    = "https://api.xero.com/oauth/AccessToken"
	baseURL           = "https://api.xero.com"
)

var client oauth.Client

// NewClient initializes the oauth Client
func NewClient(token string, key []byte) error {

	block, _ := pem.Decode(key)
	privateKey, ParseKeyErr := x509.ParsePKCS1PrivateKey(block.Bytes)

	if ParseKeyErr != nil {
		log.Printf("[xero NewClient] - Parse private key ERROR: %v", ParseKeyErr)
		return ParseKeyErr
	}

	client = oauth.Client{
		TemporaryCredentialRequestURI: requestTokenURL,
		ResourceOwnerAuthorizationURI: authorizeTokenURL,
		TokenRequestURI:               accessTokenURL,
		Header:                        http.Header{"Accept": {"application/xml"}},
		SignatureMethod:               oauth.RSASHA1,
		Credentials:                   oauth.Credentials{Token: token},
		PrivateKey:                    privateKey,
	}

	return nil
}

// PostRequest sends POST requests to xero APIs with a form as payload
func PostRequest(path string, payload string) (response string, err error) {

	form := url.Values{"xml": {payload}}

	req, reqErr := http.NewRequest("POST", baseURL, strings.NewReader(form.Encode()))
	if reqErr != nil {
		log.Printf("[xero PostRequest] - Error: %v\n", reqErr)
		return "", reqErr
	}

	req.URL.Path = path

	headerErr := client.SetAuthorizationHeader(req.Header, &client.Credentials, "POST", req.URL, nil)
	if headerErr != nil {
		log.Printf("[xero PostRequest] - SetAuthorizationHeader Error: %v\n", headerErr)
		return "", headerErr
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	resp, reqErr := client.Post(http.DefaultClient, &client.Credentials, req.URL.String(), form)
	if reqErr != nil {
		log.Printf("[xero PostRequest] - Error: %v\n", reqErr)
		return "", reqErr
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return string(body), nil
}

// Request sends requests to xero APIs
func Request(method string, path string) (response string, err error) {

	req, err := http.NewRequest(method, baseURL, nil)
	if err != nil {
		log.Printf("[xero Request] - error: %v\n", err)
		return "", err
	}

	req.URL.Opaque = path
	log.Printf("[xero Request] - URL: %s\n", req.URL.String())
	headerErr := client.SetAuthorizationHeader(req.Header, &client.Credentials, method, req.URL, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")

	if headerErr != nil {
		log.Printf("[xero Request] - SetAuthorizationHeader Error: %v\n", headerErr)
		return "", headerErr
	}

	resp, reqErr := http.DefaultClient.Do(req)
	if reqErr != nil {
		log.Printf("[xero Request] - Error: %v\n", reqErr)
		return "", reqErr
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return string(body), nil

}
