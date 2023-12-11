package report

import (
	"encoding/xml"
	"io"
	"os"
	"path"
	"strings"
)

func ParseXMLFiles(directory string) (testSuites []*TestSuites, err error) {
	dir, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for _, file := range dir {
		name := file.Name()
		if strings.HasSuffix(name, "xml") {
			testSuite, err := UnmarshalXMLFile(path.Join(directory, name))
			if err != nil {
				return nil, err
			}

			testSuites = append(testSuites, testSuite)
		}
	}
	return
}

func UnmarshalXMLFile(fileName string) (*TestSuites, error) {
	var content []byte
	testSuites := &TestSuites{}

	fd, err := os.Open(fileName)
	if err != nil {
		return testSuites, err
	}
	defer fd.Close()

	content, err = io.ReadAll(fd)
	if err != nil {
		return testSuites, err
	}

	err = xml.Unmarshal(content, testSuites)
	if err != nil {
		return testSuites, err
	}

	return testSuites, err
}
