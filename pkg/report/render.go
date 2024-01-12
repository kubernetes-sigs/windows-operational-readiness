/*
Copyright 2024 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package report

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func RenderReportTable(csv bool, testSuites []*TestSuites) string {
	var (
		err            error
		passed, failed int32
		rowHeader      = table.Row{"#", "result", "category", "time", "conformance"}
	)

	tw := table.NewWriter()
	tw.AppendHeader(rowHeader)

	for _, testSuite := range testSuites {
		var ftime string
		for _, test := range testSuite.Suites {
			if len(test.TestCases) == 0 {
				if ftime, err = parseFloat(test.Time); err != nil {
					continue
				}
				tw.AppendRow(table.Row{test.Index, "FAILED", testSuite.Category, ftime, testSuite.Name, ""})
				failed++
			}
			for _, c := range test.TestCases {
				if ftime, err = parseFloat(c.Time); err != nil {
					continue
				}
				tw.AppendRow(table.Row{test.Index, strings.ToUpper(string(c.Status)), testSuite.Category, ftime, testSuite.Name})
				if c.Status == StatusPassed {
					passed++
				}
				if c.Status == StatusFailed {
					failed++
				}
			}
		}
	}

	tw.AppendFooter(table.Row{"", fmt.Sprintf("passed: %d", passed), fmt.Sprintf("failed: %d", failed), "", ""})
	tw.SetTitle("Windows Operational Summary")
	tw.SetIndexColumn(1)
	stylePairs := []table.Style{
		table.StyleColoredBlackOnYellowWhite, table.StyleColoredYellowWhiteOnBlack,
	}
	row := make(table.Row, 2)
	for idx, style := range stylePairs {
		tw.SetStyle(style)
		tw.Style().Title.Align = text.AlignCenter
		row[idx] = tw.Render()
	}

	if csv {
		return tw.RenderCSV()
	}
	return tw.Render()
}

func parseFloat(time string) (string, error) {
	ftime, err := strconv.ParseFloat(time, 32)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%.2f", ftime), nil
}
