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
	"io/ioutil"
	"net/http"
	"strings"
)

type Client struct {
	token   string
	baseURL string
}

func NewClient(token, baseURL string) *Client {
	return &Client{token, baseURL}
}

func buildUrl(values ...string) string {
	return strings.Join(values, "/")
}

type Milestone struct {
	ID          int64  `json:"id"`
	IID         int64  `json:"iid"`
	GroupID     int64  `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	StartDate   string `json:"start_date"`
	State       string `json:"state"`
}

func (c *Client) GetMilestones(group string) ([]Milestone, error) {
	milestones := []Milestone{}

	req, err := http.NewRequest("GET", buildUrl(c.baseURL, "groups", group, "milestones"), nil)
	if err != nil {
		return milestones, err
	}

	q := req.URL.Query()
	q.Add("state", "active")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("PRIVATE-TOKEN", c.token)

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

func doRequest(req *http.Request) ([]byte, error) {
	var client = &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
