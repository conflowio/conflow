package dependency_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDependency(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dependency Suite")
}
