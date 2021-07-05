package gotau

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Time wraps the standard time.Time to allow for custom parsing from json
type Time struct {
	time.Time
}

// UnmarshalJSON implemented to allow for parsing this time object from TAU
func (t *Time) UnmarshalJSON(b []byte) error {
	primaryLayout := "2006-01-02T15:04:05.999999999-07:00"
	secondaryLayout := "2006-01-02T15:04:05.999999999-0700"
	timeAsString := strings.TrimSpace(string(b))
	timeAsString = strings.Trim(timeAsString, "\"")
	timestamp, err := time.Parse(primaryLayout, timeAsString)
	if err != nil {
		timestamp, err = time.Parse(secondaryLayout, timeAsString)
		if err != nil {
			return err
		}
	}
	t.Time = timestamp
	return nil
}

func (c *Client) apiRequest(endpoint string, params map[string][]string, body []byte, method string) ([]byte, error) {
	protocol := "http"
	if c.hasSSL {
		protocol = "https"
	}
	endpointURL := fmt.Sprintf("%s://%s:%d/api/v1/%s/", protocol, c.hostname, c.port, endpoint)
	httpClient := &http.Client{}
	buffer := bytes.NewBuffer(body)
	request, err := http.NewRequest(method, endpointURL, buffer)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", c.token))
	request.Header.Add("Content-Type", "application/json")
	_, err = request.URL.Parse(endpointURL)
	if err != nil {
		return nil, err
	}

	q := request.URL.Query()
	q.Add("format", "json")
	for key, values := range params {
		for _, item := range values {
			q.Add(key, item)
		}
	}
	request.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		body, err := ioutil.ReadAll(response.Body)
		return body, err
	}
	if response.StatusCode == 401 {
		return nil, AuthorizationError{}
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		err = GenericError{
			Err:  fmt.Sprintf("response Code %d: %s", response.StatusCode, body),
			Body: body,
			Code: response.StatusCode,
		}
		return nil, err
	}
}
