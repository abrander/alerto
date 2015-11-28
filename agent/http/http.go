package http

import (
	"net/http"
	"time"

	"github.com/abrander/alerto/agent"
)

func init() {
	agent.Register("http", NewHttp)
}

func NewHttp() agent.Agent {
	return new(Http)
}

type (
	Http struct {
		Url string `json:"url" description:"The URL to request"`
	}
)

func (h *Http) Execute(request agent.Request) agent.Result {
	start := time.Now()

	resp, err := http.Get(h.Url)
	if err != nil {
		return agent.NewResult(agent.Failed, agent.NewMeasurementCollection("time", time.Now().Sub(start)), err.Error())
	}

	c := agent.NewMeasurementCollection(
		"time", time.Now().Sub(start),
		"status", resp.StatusCode,
	)

	return agent.NewResult(agent.Ok, c, "returned %d", resp.StatusCode)
}
