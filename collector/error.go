package collector

import (
	"github.com/juju/errgo"
)

var (
	ErrSourceCodeNotFoundByFilePath = errgo.New("Source code not found by file path")

	Mask = errgo.MaskFunc()
)
