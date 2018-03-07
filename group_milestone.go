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
)

// GroupMilestone represents the Gitlab milestone
type GroupMilestone struct {
	ID          int64  `json:"id"`
	IID         int64  `json:"iid"`
	GroupID     int64  `json:"group_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	StartDate   string `json:"start_date"`
	State       string `json:"state"`
}

// MergeRequest represents the Gitlab Merge Request
type MergeRequest struct {
	ID          int64  `json:"id"`
	IID         int64  `json:"iid"`
	State       string `json:"state"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// GetMergeRequests returns an array of MilestoneMergeRequest for this GroupMilestone
func (m *GroupMilestone) GetMergeRequests(c *Client) ([]MergeRequest, error) {
	req, err := c.getRequest(buildUrl("merge_requests"))
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("state", "all")
	q.Add("milestone", m.Title)
	req.URL.RawQuery = q.Encode()

	result, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	mergeReqs := []MergeRequest{}
	err = json.Unmarshal(result, &mergeReqs)
	if err != nil {
		return nil, err
	}

	return mergeReqs, nil
}
