package pidof

import (
	"io/ioutil"
	"strings"

	"github.com/abrander/alerto/plugins"
)

func init() {
	plugins.Register("pidof", NewPidOf)
}

func NewPidOf() plugins.Plugin {
	return new(PidOf)
}

type (
	PidOf struct {
		ProcessName string `json:"processName" description:"The processname to look up"`
	}
)

func (p *PidOf) Run(transport plugins.Transport, request plugins.Request) plugins.Result {
	stdout, _, err := transport.Exec("/bin/pidof", p.ProcessName)
	if err != nil {
		return plugins.NewResult(plugins.Failed, nil, err.Error())
	}

	content, err := ioutil.ReadAll(stdout)
	if err != nil {
		return plugins.NewResult(plugins.Failed, nil, err.Error())
	}

	fields := strings.Fields(string(content))

	return plugins.NewResult(
		plugins.Ok,
		plugins.NewMeasurementCollection(
			"count",
			len(fields)),
		"%s has PID(s) %s", p.ProcessName, strings.TrimSpace(string(content)))
}

// Ensure compliance
var _ plugins.Agent = (*PidOf)(nil)
