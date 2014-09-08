package matic_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	maticPkg "github.com/zyndiecate/matic/src"
)

func TestMatic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "matic")
}

var _ = Describe("matic", func() {
	var (
		err error
	)

	BeforeEach(func() {
		// ...
	})

	AfterEach(func() {
		// ...
	})

	Describe("foo", func() {
		It("bar", func() {
			Expect(err).To(BeNil())
		})

		It("baz", func() {
			Expect(maticPkg.Foo()).To(Equal("foo"))
		})
	})
})
