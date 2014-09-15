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
			root           string
			err            error
			ctx            *collectorPkg.Ctx
			sourceCodeList []collectorPkg.SourceCode
		)

		BeforeEach(func() {
			err = nil
			root = "../fixture/simple/"

			ctx = &collectorPkg.Ctx{
				SourceCode: collectorPkg.SourceCodeCtx{
					Ext:  "go",
					Root: root,
				},
			}

			err = taskqPkg.NewQueue(ctx).RunTasks(
				taskqPkg.InSeries(
					collectorPkg.SourceCodeTask,
				),
			)

			sourceCodeList = ctx.SourceCode.SourceCodeList
		})

		Context("run source code task", func() {
			It("should not throw error", func() {
				Expect(err).To(BeNil())
			})

			It("should find 3 files", func() {
				Expect(sourceCodeList).To(HaveLen(3))
			})

			It("should find middleware/v1/middleware.go", func() {
				Expect(sourceCodeList[0].FilePath).To(Equal("../fixture/simple/middleware/v1/middleware.go"))
			})

			It("should find middleware/v1/v1.go", func() {
				Expect(sourceCodeList[1].FilePath).To(Equal("../fixture/simple/middleware/v1/v1.go"))
			})

			It("should find simple.go", func() {
				Expect(sourceCodeList[2].FilePath).To(Equal("../fixture/simple/simple.go"))
			})
		})
	})
})
