package http

import (
	"net/http"
	"time"

	"github.com/abrander/alerto/plugins"
)

func init() {
	plugins.Register("http", NewHttp)
}

func NewHttp() plugins.Plugin {
	return new(Http)
}

type (
	Http struct {
		Url string `json:"url" description:"The URL to request"`
	}
)

func (h *Http) Run(transport plugins.Transport, request plugins.Request) plugins.Result {
	start := time.Now()

	resp, err := http.Get(h.Url)
	if err != nil {
		return plugins.NewResult(plugins.Failed, plugins.NewMeasurementCollection("time", time.Now().Sub(start)), err.Error())
	}
	defer resp.Body.Close()

	c := plugins.NewMeasurementCollection(
		"time", time.Now().Sub(start),
		"status", resp.StatusCode,
	)

	return plugins.NewResult(plugins.Ok, c, "returned %d", resp.StatusCode)
}

// Ensure compliance
var _ plugins.Agent = (*Http)(nil)
