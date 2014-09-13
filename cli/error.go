package cli

import (
	"github.com/juju/errgo"
)

var (
	ErrWrongInput = errgo.New("Wrong input")

	Mask = errgo.MaskFunc()
)
