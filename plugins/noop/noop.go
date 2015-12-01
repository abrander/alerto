package noop

import (
	"time"

	"github.com/abrander/alerto/plugins"
)

func init() {
	plugins.Register("noop", NewNoop)
}

func NewNoop() plugins.Plugin {
	return new(Noop)
}

type (
	Noop struct {
		Delay time.Duration `json:"delay" description:"Amount of time to do nothing"`
	}
)

func (n *Noop) Execute(request plugins.Request) plugins.Result {
	time.Sleep(n.Delay)
	return plugins.NewResult(plugins.Ok, nil, "noop ;-)")
}
