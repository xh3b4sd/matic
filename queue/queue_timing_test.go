package queue_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TrackTiming(cb func()) int64 {
	start := time.Now().UnixNano()
	cb()
	end := time.Now().UnixNano()

	return end - start
}

func TestTimingQueue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "timing-queue")
}

var _ = Describe("timing-queue", func() {
	var (
		seriesDuration   int64
		parallelDuration int64
	)

	BeforeEach(func() {
		seriesDuration = TrackTiming(func() {
			SeriesQueue()
		})

		parallelDuration = TrackTiming(func() {
			ParallelQueue()
		})
	})

	Context("exectuting parallel queue", func() {
		It("should be faster than executing series queue", func() {
			Expect(parallelDuration).To(BeNumerically("<", seriesDuration))
		})
	})
})
