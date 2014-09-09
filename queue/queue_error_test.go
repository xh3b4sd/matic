package queue_test

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	queuePkg "github.com/zyndiecate/matic/queue"
)

func TaskErr(ctx interface{}) error {
	return fmt.Errorf("TaskErr")
}

func ErrorQueue() (*Ctx, error) {
	ctx := &Ctx{}
	q := queuePkg.NewQueue(ctx)

	err := q.RunTasks(
		q.InSeries(
			Task1,
			TaskErr,
			Task4,
			Task2,
		),
	)

	if err != nil {
		// Because we want to receive the modified context even in error cases, we
		// do not return the contexts zero value here.
		return ctx, err
	}

	return ctx, nil
}

func TestErrorQueue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "error-queue")
}

var _ = Describe("error-queue", func() {
	var (
		ctx *Ctx
		err error
	)

	BeforeEach(func() {
		ctx, err = ErrorQueue()
	})

	Context("executing error queue", func() {
		// TaskErr throws an error, thus there must be an error returned.
		It("should throw error", func() {
			Expect(err).NotTo(BeNil())
		})

		// Task1 is called before TaskErr, thus it must NOT be empty.
		It("should create correct context value for task1", func() {
			Expect(ctx.Task1).To(Equal("task1"))
		})

		// Task2 is called after TaskErr, thus it must be empty.
		It("should create no context value for task2", func() {
			Expect(ctx.Task2).To(Equal(0))
		})

		// Task3 is not called in this queue, thus it must be empty.
		It("should create no context value for task3", func() {
			Expect(ctx.Task3).To(HaveLen(0))
		})

		// Task4 is called after TaskErr, thus it must be empty.
		It("should create no context value for task4", func() {
			Expect(ctx.Task4).To(Equal(float64(0)))
		})
	})
})
