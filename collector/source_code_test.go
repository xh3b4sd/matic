package collector_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	collectorPkg "github.com/zyndiecate/matic/collector"
	taskqPkg "github.com/zyndiecate/taskq"
)

var _ = Describe("task", func() {
	Describe("source code", func() {
		var (
			wd  string
			err error
			ctx *collectorPkg.Ctx
		)

		BeforeEach(func() {
			err = nil
			wd = "../fixture/simple/"

			ctx = &collectorPkg.Ctx{
				WorkingDir: wd,
			}

			err = taskqPkg.NewQueue(ctx).RunTasks(
				taskqPkg.InSeries(
					collectorPkg.SourceCodeTask,
				),
			)
		})

		Context("run source code task", func() {
			It("should not throw error", func() {
				Expect(err).To(BeNil())
			})

			It("should find 3 files", func() {
				Expect(ctx.Files).To(HaveLen(3))
			})

			It("should find middleware/v1/middleware.go", func() {
				Expect(ctx.Files[0].Path).To(Equal("../fixture/simple/middleware/v1/middleware.go"))
			})

			It("should find middleware/v1/v1.go", func() {
				Expect(ctx.Files[1].Path).To(Equal("../fixture/simple/middleware/v1/v1.go"))
			})

			It("should find simple.go", func() {
				Expect(ctx.Files[2].Path).To(Equal("../fixture/simple/simple.go"))
			})
		})
	})
})
