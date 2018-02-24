// MIT License
//
// Copyright (c) 2018 Mario Carrion <info@carrion.ws>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Client represents the handler used for interacting with the Gitlab API
type Client struct {
	token   string
	baseURL string
}

// Error represents the value returned when there's an error returned by the
// Gitlab API
type Error struct {
	Message string `json:"message"`
}

// NewClient returns a new Client already containing the baseURL Gitlab URL
// and the access token
func NewClient(token, baseURL string) *Client {
	return &Client{token, baseURL}
}

// GetMilestones returns an array of active Milestones that belong to the
// specific group
func (c *Client) GetMilestones(group string) ([]Milestone, error) {
	milestones := []Milestone{}

	req, err := c.buildRequest(buildUrl("groups", group, "milestones"))
	if err != nil {
		return milestones, err
	}

	q := req.URL.Query()
	q.Add("state", "active")
	req.URL.RawQuery = q.Encode()

	result, err := doRequest(req)
	if err != nil {
		return milestones, err
	}

	err = json.Unmarshal(result, &milestones)
	if err != nil {
		return milestones, err
	}

	return milestones, nil
}

func (c *Client) buildRequest(url string) (*http.Request, error) {
	r, err := http.NewRequest("GET", buildUrl(c.baseURL, url), nil)
	if err == nil {
		r.Header.Add("PRIVATE-TOKEN", c.token)
	}

	return r, err
}

func buildUrl(values ...string) string {
	return strings.Join(values, "/")
}

func doRequest(req *http.Request) ([]byte, error) {
	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return b, err
	}

	if resp.StatusCode != 200 {
		return []byte{}, createErrorFromResponse(b)
	}
	return b, nil
}

func createErrorFromResponse(r []byte) error {
	errStruct := Error{}
	err := json.Unmarshal(r, &errStruct)
	if err != nil {
		return err
	}

	return fmt.Errorf(errStruct.Message)
}
