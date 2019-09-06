// Package utility contains methods to make HTTP request
package utility

import (
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

//---------------------------------------------------------------------------------------

var httpClient *http.Client

//---------------------------------------------------------------------------------------

func init() {
	httpClient = &http.Client{}
} // init

//---------------------------------------------------------------------------------------

var (
	requestConstructor = Request
)

//---------------------------------------------------------------------------------------

// Request constructs http.Request instance
func Request(url, method string, requestBody io.Reader, headers *http.Header) (*http.Request, error) {
	// Create Request
	req, err := http.NewRequest(method, url, requestBody)
	if err != nil {
		log.Errorf("Error in creating request object: %v", err)
		return nil, err
	}
	if headers != nil {
		req.Header = *headers
	}
	return req, nil
} // Request

//---------------------------------------------------------------------------------------

// SendRequest sends Http Request at given url with given method, request body and headers
func SendRequest(url, method string, requestBody io.Reader, headers *http.Header) (*[]byte, int, error) {
	// Create Request
	req, err := requestConstructor(url, method, requestBody, headers)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// Fetch Request
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Errorf("Error in sending request to %s : %v", url, err)
		var statusCode = http.StatusBadRequest
		return nil, statusCode, err
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body. URL: %s. Error: %v", url, err)
		return nil, http.StatusInternalServerError, err
	}
	return &respBody, resp.StatusCode, nil
} // sendRequest
