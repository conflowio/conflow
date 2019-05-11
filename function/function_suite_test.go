package function_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFunction(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Function Suite")
}
