package collector

type ServerNameCtx struct {
	// Server name of the created middleware server.
	Name     string
	FilePath string
}

func ServerNameTask(ctx interface{}) error {
	Verbosef("Searching server name")

	// each source code of package import path
	for _, packageImport := range ctx.(*Ctx).PackageImport.PackageImportList {
		astTree, err := astTreeByFile(packageImport.FilePath, ctx.(*Ctx).SourceCode.SourceCodeList)
		if err != nil {
			return Mask(err)
		}

		// NOTICE: here we assume there is only one created middleware server.
		// Maybe that is not true for all cases.
		srvName := serverName(packageImport.PkgName, astTree)
		if srvName != "" {
			ctx.(*Ctx).ServerName.Name = srvName
			ctx.(*Ctx).ServerName.FilePath = packageImport.FilePath
			break
		}
	}

	return nil
}
