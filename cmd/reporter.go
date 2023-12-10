/*
Copyright 2023 The Kubernetes Authors.

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

package cmd

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"sigs.k8s.io/windows-operational-readiness/pkg/report"
)

var (
	directory string
	failed    int
	passed    int
)

var reporterCmd = &cobra.Command{
	Use:   "reporter",
	Short: "Parse XML Junit results and print a columnar output",
	Long:  `Parse XML Junit results and print a columnar output`,
	Run: func(cmd *cobra.Command, args []string) {
		suites, err := report.ParseXMLFiles(directory)
		if err != nil {
			log.Fatal(err)
		}
		for _, n := range suites {
			presentTestSuite(n)
		}
		fmt.Printf("Passed tests: %d\n", passed)
		fmt.Printf("Failed tests: %d\n", failed)
	},
}

func init() {
	reporterCmd.Flags().StringVar(&directory, "dir", "", "directory to search in for xml files")
}

func presentTestSuite(r *report.TestSuites) {
	for _, t := range r.Suites {
		if len(t.TestCases) == 0 {
			ftime, err := strconv.ParseFloat(t.Time, 32)
			if err != nil {
				continue
			}
			fmt.Printf("[FAILED] | %s - %s | (%2.2fs) | %s\n", t.Index, r.Category, ftime, r.Name)
			failed++
		}

		for _, c := range t.TestCases {
			ftime, err := strconv.ParseFloat(c.Time, 32)
			if err != nil {
				continue
			}
			fmt.Printf("[%s] | %s - %s | (%2.2fs) | %s | %s\n", strings.ToUpper(string(c.Status)), t.Index, r.Category, ftime, r.Name, c.Name)
			if c.Status == report.StatusPassed {
				passed++
			}
			if c.Status == report.StatusFailed {
				failed++
			}
		}
	}
}
