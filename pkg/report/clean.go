package report

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"go.uber.org/zap"
)

// CleanupJUnitXML removes unnecessary skipped tests from the report.
func CleanupJUnitXML(path, category, testName string, index int) error {
	zap.L().Info("Getting XML content from file", zap.String("path", path), zap.String("category", category))

	content, err := getXMLContent(path)
	if err != nil {
		return err
	}

	var testSuites TestSuites
	if err = xml.Unmarshal(content, &testSuites); err != nil {
		return err
	}
	for i, suite := range testSuites.Suites {
		var cleaned []TestCase
		for _, test := range suite.TestCases {
			if test.Status != StatusSkipped &&
				test.Name != "[SynchronizedBeforeSuite]" &&
				test.Name != "[SynchronizedAfterSuite]" &&
				test.Name != "[ReportBeforeSuite]" &&
				test.Name != "[ReportAfterSuite] Kubernetes e2e suite report" {
				cleaned = append(cleaned, test)
			}
		}

		zap.L().Info("Saving cleaned tests.", zap.Int("number", len(cleaned)), zap.Int("index", index))
		testSuites.Suites[i].Index = fmt.Sprint(index)
		testSuites.Suites[i].TestCases = cleaned
	}

	// save the category metadata
	testSuites.Category = category
	testSuites.Name = testName

	// write back the cleaned up YAML to a writer
	var cleanContent []byte
	if cleanContent, err = xml.MarshalIndent(testSuites, "  ", "    "); err != nil {
		return err
	}

	return writeFileContent(path, cleanContent)
}

// writeFileContent save the content to a file.
func writeFileContent(path string, content []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = file.Write(content); err != nil {
		return err
	}
	return nil
}

// getXMLContent returns the content in bytes of an existent file.
func getXMLContent(path string) ([]byte, error) {
	var err error
	// check if file exists
	if _, err = os.Stat(path); os.IsNotExist(err) {
		return []byte{}, err
	}

	var file *os.File
	if file, err = os.Open(path); err != nil {
		return []byte{}, err
	}
	defer file.Close()

	var content []byte
	if content, err = io.ReadAll(file); err != nil {
		return []byte{}, err
	}

	return content, nil
}
