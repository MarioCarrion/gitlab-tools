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
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MarioCarrion/gitlab-tools"
)

type params struct {
	token   string
	baseURL string
	group   string
	from    string
	to      string
	dueDate string
}

func main() {
	params := parseParams()

	c := gitlab.NewClient(params.token, params.baseURL)
	milestones, err := c.GetMilestones(params.group)
	if err != nil {
		log.Fatalln("error", err)
	}

	var toMilestone, fromMilestone gitlab.GroupMilestone

	for _, milestone := range milestones {
		if params.to == milestone.Title {
			toMilestone = milestone
		}
		if params.from == milestone.Title {
			fromMilestone = milestone
		}

	}

	if fromMilestone.ID == 0 {
		log.Fatalf("milestone '%s' not found, exiting\n", params.from)
	}

	if toMilestone.ID == 0 {
		fmt.Printf("milestone '%s' not found, creating it\n", params.to)
		toMilestone.StartDate = fromMilestone.DueDate
		toMilestone.DueDate = params.dueDate
		toMilestone.Title = params.to
		toMilestone, err = c.CreateMilestone(fromMilestone.GroupID, toMilestone)
		if err != nil {
			log.Fatalln("error creating new milestone", err)
		}
	}

	b, err := json.MarshalIndent(toMilestone, "", "  ")
	fmt.Printf("to:\n%s\n", string(b))

	b, err = json.MarshalIndent(fromMilestone, "", "  ")
	fmt.Printf("from:\n%s\n", string(b))

	//

	issues, err := fromMilestone.GetMergeRequests(c)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	} else {
		b, err = json.MarshalIndent(issues, "", "  ")
		fmt.Printf("%s\n", string(b))
	}
}

func parseParams() params {
	r := params{}

	flag.StringVar(&r.token, "token", "", "private Gitlab token")
	flag.StringVar(&r.baseURL, "baseURL", "", "base Gitlab URL")
	flag.StringVar(&r.group, "group", "", "Gitlab group name")
	flag.StringVar(&r.to, "to", "", "milestone to move merge requests to")
	flag.StringVar(&r.from, "from", "", "milestone to move merge requests from")
	flag.StringVar(&r.dueDate, "due", "", "due date for the new milestone")

	flag.Parse()

	retToken := r.token
	retBaseURL := r.baseURL

	if len(retToken) == 0 {
		r.token = os.Getenv("GITLAB_TOOLS_TOKEN")
	}
	if len(retBaseURL) == 0 {
		r.baseURL = os.Getenv("GITLAB_TOOLS_BASE_URL")
	}

	fatalWhenInvalidString("token", r.token)
	fatalWhenInvalidString("baseURL", r.baseURL)
	fatalWhenInvalidString("group", r.group)
	fatalWhenInvalidString("to", r.to)
	fatalWhenInvalidString("from", r.from)
	fatalWhenInvalidString("due", r.dueDate)

	return r
}

func fatalWhenInvalidString(n, s string) {
	if len(s) == 0 {
		log.Fatalf("Required parameter missing: '%s'\n", n)
	}
}
