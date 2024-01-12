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

	"github.com/spf13/cobra"
	"sigs.k8s.io/windows-operational-readiness/pkg/report"
)

var (
	directory string
	csv       bool
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
		fmt.Println(report.RenderReportTable(csv, suites))
	},
}

func init() {
	reporterCmd.Flags().StringVar(&directory, "dir", "", "directory to search in for xml files.")
	reporterCmd.Flags().BoolVar(&csv, "csv", false, "return the table as CSV.")
}
