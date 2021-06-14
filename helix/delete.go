package helix

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// DeleteRequest runs a raw GET request against the twitch helix pass thru, and returns true if the delete worked.
func (c *Client) DeleteRequest(endpoint string, params map[string][]string) (bool, error) {
	protocol := "http"
	if c.hasSSL {
		protocol = "https"
	}
	endpointURL := fmt.Sprintf("%s://%s:%d/api/twitch/helix/%s/", protocol, c.hostname, c.port, endpoint)
	httpClient := &http.Client{}
	request, err := http.NewRequest("DELETE", endpointURL, nil)
	if err != nil {
		return false, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("Token %s", c.token))
	_, err = request.URL.Parse(endpointURL)
	if err != nil {
		return false, err
	}

	q := request.URL.Query()
	for key, values := range params {
		for _, item := range values {
			q.Add(key, item)
		}
	}
	request.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(request)
	if err != nil {
		return false, err
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return true, nil
	}
	if response.StatusCode == 401 {
		return false, AuthorizationError{}
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		err = GenericError{
			err:  fmt.Sprintf("response code %d: %s", response.StatusCode, body),
			body: body,
			code: response.StatusCode,
		}
		return false, err
	}
}

func (c *Client) DeleteCustomReward(broadcasterID, ID string) (bool, error) {
	broadcasterID = strings.TrimSpace(broadcasterID)
	ID = strings.TrimSpace(ID)
	if broadcasterID == "" {
		return false, BadRequestError{
			"invalid request, broadcaster can't be blank",
		}
	}
	if ID == "" {
		return false, BadRequestError{
			"invalid request, id can't be blank",
		}
	}

	params := map[string][]string{
		"broadcaster_id": {broadcasterID},
		"id":             {ID},
	}

	return c.DeleteRequest("channel_points/custom_rewards", params)

}
