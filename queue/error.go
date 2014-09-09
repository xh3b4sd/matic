package queue

import (
	"github.com/juju/errgo"
)

var (
	ErrNoTasks = errgo.New("No tasks")

	Mask = errgo.MaskFunc()
)
