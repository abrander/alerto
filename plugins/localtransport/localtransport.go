package localtransport

import (
	"bytes"
	"io"
	"os/exec"

	"github.com/abrander/alerto/plugins"
)

func init() {
	plugins.Register("localtransport", NewLocalTransport)
}

func NewLocalTransport() plugins.Plugin {
	return new(LocalTransport)
}

type (
	LocalTransport struct {
	}
)

func (l *LocalTransport) Exec(cmd string, arguments ...string) (io.Reader, io.Reader, error) {
	command := exec.Command(cmd, arguments...)

	var out bytes.Buffer
	command.Stdout = &out

	stderr, err := command.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	err = command.Run()

	return &out, stderr, err
}

// Ensure compliance
var _ plugins.Transport = (*LocalTransport)(nil)
