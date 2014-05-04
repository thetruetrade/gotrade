package gotrade_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGotrade(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gotrade Suite")
}
