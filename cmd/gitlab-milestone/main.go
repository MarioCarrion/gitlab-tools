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

	"github.com/MarioCarrion/gitlab-tools"
)

func main() {
	token := flag.String("token", "", "private Gitlab token")
	baseURL := flag.String("baseURL", "", "base Gitlab URL")
	group := flag.String("group", "", "Gitlab group name")

	flag.Parse()

	if len(*token) == 0 {
		log.Fatalln("Required parameter missing: 'token'")
	}
	if len(*baseURL) == 0 {
		log.Fatalln("Required parameter missing: 'baseURL'")
	}
	if len(*group) == 0 {
		log.Fatalln("Required parameter missing: 'group'")
	}

	git := gitlab.NewClient(*token, *baseURL)
	milestones, err := git.GetMilestones(*group)
	if err != nil {
		log.Fatalln("error ", err)
	}

	for _, milestone := range milestones {
		fmt.Printf("%+v\n", milestone)
	}
}
