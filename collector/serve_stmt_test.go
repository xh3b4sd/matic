package collector_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	collectorPkg "github.com/zyndiecate/matic/collector"
	taskqPkg "github.com/zyndiecate/taskq"
)

var _ = Describe("serve-stmt", func() {
	Describe("source code", func() {
		var (
			ctx *collectorPkg.Ctx
			wd  string
			err error
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
					collectorPkg.PackageImportTask,
					collectorPkg.ServerNameTask,
					collectorPkg.ServeStmtTask,
				),
			)
		})

		Context("run source code task", func() {
			It("should not throw error", func() {
				Expect(err).To(BeNil())
			})

			It("should find 3 serve statements", func() {
				Expect(ctx.Files[0].ServeStmts).To(HaveLen(0))
				Expect(ctx.Files[1].ServeStmts).To(HaveLen(2))
				Expect(ctx.Files[2].ServeStmts).To(HaveLen(1))
			})

			It("should find serve statement of 'hello' route", func() {
				Expect(ctx.Files[1].ServeStmts[0].Method).To(Equal("GET"))
				Expect(ctx.Files[1].ServeStmts[0].Path).To(Equal("/v1/hello"))
				Expect(ctx.Files[1].ServeStmts[0].Middlewares).To(HaveLen(1))
				Expect(ctx.Files[1].ServeStmts[0].Middlewares[0].Type).To(Equal("V1"))
				Expect(ctx.Files[1].ServeStmts[0].Middlewares[0].Name).To(Equal("Hello"))
			})

			It("should find serve statement of 'world' route", func() {
				Expect(ctx.Files[1].ServeStmts[1].Method).To(Equal("GET"))
				Expect(ctx.Files[1].ServeStmts[1].Path).To(Equal("/v1/world"))
				Expect(ctx.Files[1].ServeStmts[1].Middlewares).To(HaveLen(2))
				Expect(ctx.Files[1].ServeStmts[1].Middlewares[0].Type).To(Equal(""))
				Expect(ctx.Files[1].ServeStmts[1].Middlewares[0].Pkg).To(Equal("v1"))
				Expect(ctx.Files[1].ServeStmts[1].Middlewares[0].Name).To(Equal("Foo"))
				Expect(ctx.Files[1].ServeStmts[1].Middlewares[1].Type).To(Equal("V1"))
				Expect(ctx.Files[1].ServeStmts[1].Middlewares[1].Pkg).To(Equal(""))
				Expect(ctx.Files[1].ServeStmts[1].Middlewares[1].Name).To(Equal("World"))
			})

			It("should find serve statement of 'hello world' route", func() {
				Expect(ctx.Files[2].ServeStmts[0].Method).To(Equal("GET"))
				Expect(ctx.Files[2].ServeStmts[0].Path).To(Equal("/v1/hello-world"))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares).To(HaveLen(4))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares[0].Type).To(Equal("V1"))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares[0].Pkg).To(Equal(""))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares[0].Name).To(Equal("HelloWorldOne"))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares[1].Type).To(Equal("V1"))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares[1].Pkg).To(Equal(""))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares[1].Name).To(Equal("HelloWorldTwo"))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares[2].Type).To(Equal(""))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares[2].Pkg).To(Equal("v1Pkg"))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares[2].Name).To(Equal("Foo"))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares[3].Type).To(Equal("V1"))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares[3].Pkg).To(Equal(""))
				Expect(ctx.Files[2].ServeStmts[0].Middlewares[3].Name).To(Equal("HelloWorldThree"))
			})
		})
	})
})
