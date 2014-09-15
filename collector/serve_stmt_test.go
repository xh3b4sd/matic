package collector_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	collectorPkg "github.com/zyndiecate/matic/collector"
	taskqPkg "github.com/zyndiecate/taskq"
)

var _ = Describe("serve-stmt", func() {
	Describe("source code", func() {
		var (
			ctx           *collectorPkg.Ctx
			root          string
			err           error
			serveStmtList []collectorPkg.ServeStmt
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
					collectorPkg.PackageImportTask,
					collectorPkg.ServerNameTask,
					collectorPkg.ServeStmtTask,
				),
			)

			serveStmtList = ctx.ServeStmt.ServeStmtList
		})

		Context("run source code task", func() {
			It("should not throw error", func() {
				Expect(err).To(BeNil())
			})

			It("should find 3 serve statements", func() {
				Expect(ctx.ServeStmt.ServeStmtList).To(HaveLen(3))
			})

			It("should find serve statement of 'hello' route", func() {
				Expect(serveStmtList[0].Method).To(Equal("GET"))
				Expect(serveStmtList[0].Path).To(Equal("/v1/hello"))
				Expect(serveStmtList[0].Middlewares).To(HaveLen(1))
				Expect(serveStmtList[0].Middlewares[0].FuncSel).To(Equal("MiddlewareOne"))
			})

			It("should find serve statement of 'world' route", func() {
				Expect(serveStmtList[1].Method).To(Equal("GET"))
				Expect(serveStmtList[1].Path).To(Equal("/v1/world"))
				Expect(serveStmtList[1].Middlewares).To(HaveLen(1))
				Expect(serveStmtList[1].Middlewares[0].FuncSel).To(Equal("MiddlewareTwo"))
			})

			It("should find serve statement of 'hello world' route", func() {
				fmt.Printf("%#v\n", ctx.ServeStmt.ServeStmtList)
				Expect(serveStmtList[2].Method).To(Equal("GET"))
				Expect(serveStmtList[2].Path).To(Equal("/v1/hello-world"))
				Expect(serveStmtList[2].Middlewares).To(HaveLen(2))
				Expect(serveStmtList[2].Middlewares[0].FuncSel).To(Equal("MiddlewareOne"))
				Expect(serveStmtList[2].Middlewares[1].FuncSel).To(Equal("MiddlewareTwo"))
			})
		})
	})
})
