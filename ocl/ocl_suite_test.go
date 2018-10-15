package ocl_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOcl(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ocl Suite")
}
