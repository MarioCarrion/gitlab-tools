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
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MarioCarrion/gitlab-tools"
)

func main() {
	token, baseURL, group := parseArgs()

	c := gitlab.NewClient(token, baseURL)
	milestones, err := c.GetMilestones(group)
	if err != nil {
		log.Fatalln("error", err)
	}

	for _, milestone := range milestones {
		fmt.Printf("%+v\n", milestone)

		issues, err := milestone.GetMergeRequests(c)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			continue
		}
		fmt.Printf("%+v\n", issues)
	}
}

func parseArgs() (string, string, string) {
	token := flag.String("token", "", "private Gitlab token")
	baseURL := flag.String("baseURL", "", "base Gitlab URL")
	group := flag.String("group", "", "Gitlab group name")

	flag.Parse()

	retToken := *token
	retBaseURL := *baseURL

	if len(retToken) == 0 {
		retToken = os.Getenv("GITLAB_TOOLS_TOKEN")
	}
	if len(retBaseURL) == 0 {
		retBaseURL = os.Getenv("GITLAB_TOOLS_BASE_URL")
	}

	fatalWhenInvalidString("token", retToken)
	fatalWhenInvalidString("baseURL", retBaseURL)
	fatalWhenInvalidString("group", *group)

	return retToken, retBaseURL, *group
}

func fatalWhenInvalidString(n, s string) {
	if len(s) == 0 {
		log.Fatalf("Required parameter missing: '%s'\n", n)
	}
}
