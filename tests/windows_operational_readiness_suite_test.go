package main_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2" //nolint
	. "github.com/onsi/gomega"    //nolint
	_ "sigs.k8s.io/windows-operational-readiness/tests/testcases"
	_ "sigs.k8s.io/windows-operational-readiness/tests/xml"
)

func TestWindowsOperationalReadiness(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "WindowsOperationalReadiness Suite")
}
