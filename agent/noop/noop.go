package noop

import (
	"github.com/abrander/alerto/agent"
)

func init() {
	agent.Register("noop", NewNoop)
}

func NewNoop() agent.Agent {
	return new(Noop)
}

type (
	Noop struct{}
)

func (n *Noop) Execute(request agent.Request) agent.Result {
	return agent.NewResult(agent.Ok, nil, "noop ;-)")
}
