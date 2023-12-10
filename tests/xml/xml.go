package xml

import (
	. "github.com/onsi/ginkgo/v2" //nolint
	. "github.com/onsi/gomega"    //nolint
	"sigs.k8s.io/windows-operational-readiness/pkg/report"
)

var _ = Describe("Xml", func() {

	Describe("Must be parsed from junit", func() {
		Context("when passed", func() {
			It("must have one or more testcases", func() {
				var file = "./samples/junit_1101_valid.xml"
				r, err := report.UnmarshalXMLFile(file)
				Expect(err).To(BeNil())

				Expect(r.Category).To(Equal("Core.Concurrent"))

				for _, suite := range r.Suites {
					Expect(suite.Name).To(Equal("Kubernetes e2e suite"))
					Expect(suite.Tests).To(Equal("7391"))
					Expect(suite.Time).To(Equal("22.44"))
					Expect(len(suite.TestCases)).To(BeNumerically(">=", 1))
				}
			})
		})

		Context("when failed", func() {
			It("must render the error", func() {
				var file = "./samples/junit_1101_fail.xml"
				r, err := report.UnmarshalXMLFile(file)
				Expect(err).To(BeNil())

				Expect(r.Category).To(Equal("Core.Concurrent"))
				Expect(r.Failures).To(Equal("1"))

				for _, suite := range r.Suites {
					Expect(suite.Name).To(Equal("Kubernetes e2e suite"))
					Expect(suite.Tests).To(Equal("7391"))
					Expect(suite.Time).To(Equal("34.75"))
					Expect(len(suite.Failures)).To(Equal(1))
					Expect(len(suite.TestCases)).To(BeNumerically(">=", 1))
					for _, test := range suite.TestCases {
						Expect(len(test.Failure)).To(Equal(1))
						Expect(test.Status).To(Equal(report.StatusFailed))
					}
				}
			})
		})
	})

})
