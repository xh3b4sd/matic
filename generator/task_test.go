package generator_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	generatorPkg "github.com/zyndiecate/matic/generator"
	taskqPkg "github.com/zyndiecate/taskq"
)

func TestTask(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "task")
}

var _ = Describe("task", func() {
	Describe("source code", func() {
		var (
			ctx  *generatorPkg.Ctx
			root string
			err  error
		)

		BeforeEach(func() {
			err = nil
			root = "../fixture/simple/"

			ctx = &generatorPkg.Ctx{
				SourceCodeCtx: generatorPkg.SourceCodeCtx{
					Root: root,
				},
			}

			err = taskqPkg.NewQueue(ctx).RunTasks(
				taskqPkg.InSeries(
					generatorPkg.SourceCodeTask,
				),
			)
		})

		Context("run source code task", func() {
			It("should not throw error", func() {
				Expect(err).To(BeNil())
			})

			It("should find 3 files", func() {
				Expect(ctx.SourceCodeCtx.SourceCodeList).To(HaveLen(3))
			})

			It("should find middleware/v1/middleware.go", func() {
				Expect(ctx.SourceCodeCtx.SourceCodeList[0].Path).To(Equal("../fixture/simple/middleware/v1/middleware.go"))
			})

			It("should find middleware/v1/v1.go", func() {
				Expect(ctx.SourceCodeCtx.SourceCodeList[1].Path).To(Equal("../fixture/simple/middleware/v1/v1.go"))
			})

			It("should find simple.go", func() {
				Expect(ctx.SourceCodeCtx.SourceCodeList[2].Path).To(Equal("../fixture/simple/simple.go"))
			})
		})
	})
})
