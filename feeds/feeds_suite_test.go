package feeds_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFeeds(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Feeds Suite")
}
