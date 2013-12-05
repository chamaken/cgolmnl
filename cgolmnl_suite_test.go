package cgolmnl_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestCgolmnl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cgolmnl Suite")
}
