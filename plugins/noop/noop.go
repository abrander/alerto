package noop

import (
	"bytes"
	"io"
	"net"
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

func (n *Noop) Run(transport plugins.Transport, request plugins.Request) plugins.Result {
	time.Sleep(n.Delay)
	return plugins.NewResult(plugins.Ok, nil, "noop ;-)")
}

func (n *Noop) Exec(cmd string, arguments ...string) (io.Reader, io.Reader, error) {
	var stdoutBuf, stderrBuf bytes.Buffer
	return &stdoutBuf, &stderrBuf, nil
}

func (n *Noop) Dial(network string, address string) (net.Conn, error) {
	return nil, nil
}

func (n *Noop) ReadFile(path string) (io.Reader, error) {
	return nil, nil
}

// Ensure compliance
var _ plugins.Agent = (*Noop)(nil)
var _ plugins.Transport = (*Noop)(nil)
