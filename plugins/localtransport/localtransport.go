package localtransport

import (
	"bytes"
	"io"
	"net"
	"os/exec"
	"time"

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

func (l *LocalTransport) Dial(network string, address string) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	return dialer.Dial(network, address)
}

// Ensure compliance
var _ plugins.Transport = (*LocalTransport)(nil)
