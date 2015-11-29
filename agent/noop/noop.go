package noop

import (
	"time"

	"github.com/abrander/alerto/agent"
)

func init() {
	agent.Register("noop", NewNoop)
}

func NewNoop() agent.Agent {
	return new(Noop)
}

type (
	Noop struct {
		Delay time.Duration `json:"delay" description:"Amount of time to do nothing"`
	}
)

func (n *Noop) Execute(request agent.Request) agent.Result {
	time.Sleep(n.Delay)
	return agent.NewResult(agent.Ok, nil, "noop ;-)")
}
