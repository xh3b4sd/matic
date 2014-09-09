package queue_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	queuePkg "github.com/zyndiecate/matic/queue"
)

func ParallelQueue() (*Ctx, error) {
	ctx := &Ctx{}
	q := queuePkg.NewQueue(ctx)

	err := q.RunTasks(
		q.InParallel(
			Task1,
			Task3,
			Task4,
			Task2,
		),
	)

	if err != nil {
		return &Ctx{}, err
	}

	return ctx, nil
}

func TestParallelQueue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "parallel-queue")
}

var _ = Describe("parallel-queue", func() {
	var (
		ctx *Ctx
		err error
	)

	BeforeEach(func() {
		ctx, err = ParallelQueue()
	})

	Context("executing mixed queue", func() {
		It("should not throw error", func() {
			Expect(err).To(BeNil())
		})

		It("should create correct context value for task1", func() {
			Expect(ctx.Task1).To(Equal("task1"))
		})

		It("should create correct context value for task2", func() {
			Expect(ctx.Task2).To(Equal(2))
		})

		It("should create correct context value for task3", func() {
			Expect(ctx.Task3).To(Equal([]string{"task3"}))
		})

		It("should create correct context value for task4", func() {
			Expect(ctx.Task4).To(Equal(4.4))
		})
	})
})
