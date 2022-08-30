package report

import (
	"encoding/xml"
)

type Status string

const (
	StatusPassed  Status = "passed"
	StatusSkipped Status = "skipped"
	StatusFailed  Status = "failed"
	StatusError   Status = "error"
)

type TestCounter struct {
	Tests    string `xml:"tests,attr"`
	Disabled string `xml:"disabled,attr"`
	Errors   string `xml:"errors,attr"`
	Failures string `xml:"failures,attr"`
	Time     string `xml:"time,attr"`
}

type TestSuites struct {
	XMLName xml.Name `xml:"testsuites"`
	Suites  []Suite  `xml:"testsuite"`
	TestCounter
}

// Suite represents a logical grouping (suite) of tests.
type Suite struct {
	XMLName   xml.Name   `xml:"testsuite"`
	Name      string     `xml:"name,attr"`
	Package   string     `xml:"package,attr"`
	TestCases []TestCase `xml:"testcase,omitempty"`
	TestCounter
}

// TestCase represents the results of a single test run.
type TestCase struct {
	XMLName xml.Name `xml:"testcase"`

	Name      string `xml:"name,attr"`
	Classname string `xml:"classname,attr"`
	Time      string `xml:"time,attr"`
	Status    Status `xml:"status,attr"`

	Failure   []Failure `xml:"failure,omitempty"`
	SystemErr string    `xml:"system-err,omitempty"`
}

type Failure struct {
	Message string `xml:"message,attr,omitempty"`
	Type    string `xml:"type,attr,omitempty"`
}
